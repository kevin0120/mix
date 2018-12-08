package minio

import (
	"encoding/json"
	"fmt"
	"github.com/masami10/rush/services/storage"
	"github.com/minio/minio-go"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	bucket      string

	DB         *storage.Service
	minio      *minio.Client
	saveBuffer chan *ControllerCurve
	closing    chan struct{}
	wg         sync.WaitGroup
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag:       d,
		saveBuffer: make(chan *ControllerCurve, 1024),
	}
	s.configValue.Store(c)
	s.bucket = c.Bucket
	return s

}

func (s *Service) Open() error {
	c := s.config()
	client, err := minio.New(c.URL, c.Access, c.Secret, c.Secure)
	if err != nil {
		return fmt.Errorf("create minio fail %s", err.Error())
	}
	s.minio = client

	go s.saveProcess()

	// 启动重传服务
	go s.TaskReupload()

	return nil
}

func (s *Service) Close() error {
	s.closing <- struct{}{}
	s.wg.Wait()

	return nil
}

func (s *Service) Save(curve *ControllerCurve) {
	s.saveBuffer <- curve
}

// 异步保存
func (s *Service) saveProcess() {
	for {
		select {
		case data := <-s.saveBuffer:
			s.handleSave(data)

		case <-s.closing:
			s.wg.Done()
			return
		}
	}
}

// 处理保存
func (s *Service) handleSave(curve *ControllerCurve) {
	// 保存对象存储
	content, _ := json.Marshal(curve.CurveContent)
	str_content := string(content)
	err := s.Upload(curve.CurveFile, str_content)

	// 保存本地数据库
	has_upload := true
	if err != nil {
		has_upload = false
		s.diag.Error("上传曲线失败", err)
	} else {
		s.diag.Debug("上传曲线成功")
	}

	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation("2006-01-02 15:04:05", curve.UpdateTime, loc)

	//utc, _ := time.LoadLocation("")

	dbCurve := storage.Curves{
		ResultID:   curve.ResultID,
		Count:      curve.Count,
		CurveFile:  curve.CurveFile,
		CurveData:  str_content,
		HasUpload:  has_upload,
		UpdateTime: dt.UTC(),
	}

	err = s.DB.Store(dbCurve)
	if err != nil {
		s.diag.Error("缓存曲线失败", err)
	} else {
		s.diag.Debug("缓存曲线成功")
	}

}

func (s *Service) Upload(obj string, data string) error {
	isExist, err := s.minio.BucketExists(s.bucket)
	if err != nil || !isExist {
		return fmt.Errorf("Bucket %s not exist err msg: %s ", s.bucket, err.Error())
	}
	reader := strings.NewReader(data)
	_, e := s.minio.PutObject(s.bucket, obj, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/json"})
	if e != nil {
		return fmt.Errorf("Put Object %s fail ", obj)
	}
	return nil
}

func (s *Service) TaskReupload() {
	for {

		curves, err := s.DB.ListUnuploadCurves()
		if err == nil {
			for _, v := range curves {
				err = s.Upload(v.CurveFile, v.CurveData)
				if err != nil {
					s.diag.Error(fmt.Sprintf("curve reupload failed, curve_id:%d result_id:%d", v.Id, v.ResultID), err)
				} else {
					v.HasUpload = true
					s.DB.UpdateCurve(&v)
				}
			}
		}

		time.Sleep(time.Duration(s.config().ReuploadItv))
	}
}
