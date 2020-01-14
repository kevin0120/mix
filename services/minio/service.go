package minio

import (
	"fmt"
	"github.com/minio/minio-go"
	"go.uber.org/atomic"
	"strings"
	"time"
)

type Diagnostic interface {
	Error(msg string, err error)
	Debug(msg string)
	Info(msg string)
}

type Service struct {
	configValue    atomic.Value
	diag           Diagnostic
	bucket         string
	storageService IStorageService
	minio          *minio.Client
	closing        chan struct{}
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func NewService(c Config, d Diagnostic, db IStorageService) *Service {
	s := &Service{
		diag:           d,
		closing:        make(chan struct{}, 1),
		minio:          nil,
		storageService: db,
	}
	s.configValue.Store(c)

	s.bucket = c.Bucket
	return s

}

func (s *Service) Open() error {
	err := s.initClient()
	if err != nil {
		return err
	}

	// 启动重传服务
	go s.TaskReupload()

	return nil
}

func (s *Service) Close() error {
	s.closing <- struct{}{}
	return nil
}

func (s *Service) initClient() error {
	c := s.config()
	client, err := minio.New(c.URL, c.Access, c.Secret, c.Secure)
	if err != nil {
		return fmt.Errorf("Create Minio Client Failed: %s ", err.Error())
	}

	s.minio = client

	return nil
}

func (s *Service) upload(obj string, data string) error {
	isExist, err := s.minio.BucketExists(s.bucket)
	if err != nil {
		return fmt.Errorf("Get Bucket Failed: %s ", err.Error())
	}

	if !isExist {
		return fmt.Errorf("Bucket %s Not Exist ", s.bucket)
	}

	reader := strings.NewReader(data)
	_, e := s.minio.PutObject(s.bucket, obj, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/json"})
	if e != nil {
		return fmt.Errorf("Put Object %s Failed ", obj)
	}
	return nil
}

func (s *Service) doReupload() {
	curves, err := s.storageService.ListUnuploadCurves()
	if err != nil {
		s.diag.Debug(fmt.Sprintf("Get Curve Failed: %s", err.Error()))
	}

	for _, v := range curves {
		err = s.upload(v.CurveFile, v.CurveData)
		if err != nil {
			s.diag.Error(fmt.Sprintf("上传曲线失败 工具:%s 对应拧紧ID:%s", v.ToolSN, v.TighteningID), err)
		} else {
			v.HasUpload = true
			if _, err := s.storageService.UpdateCurve(&v); err != nil {
				s.diag.Error(fmt.Sprintf("本地更新曲线失败 工具:%s 对应拧紧ID:%s", v.ToolSN, v.TighteningID), err)
			} else {
				s.diag.Info(fmt.Sprintf("上传曲线成功 工具:%s 对应拧紧ID:%s", v.ToolSN, v.TighteningID))
			}
		}
	}
}

func (s *Service) TaskReupload() {
	for {
		select {
		case <-time.After(time.Duration(s.config().ReuploadItv)):
			s.doReupload()

		case <-s.closing:
			return
		}
	}
}
