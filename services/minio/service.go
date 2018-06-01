package minio

import (
	"sync/atomic"
	"github.com/minio/minio-go"
	"fmt"
	"strings"
)

type Diagnostic interface {
	Error(msg string, err error)
}

type Service struct {
	configValue atomic.Value
	diag 		Diagnostic
	bucket 		string

	minio 		*minio.Client
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag: d,
	}
	s.configValue.Store(c)
	s.bucket = c.Bucket
	return s

}

func (s *Service) Open () error {
	c := s.config()
	client, err := minio.New(c.URL, c.Access, c.Secret, c.Secure)
	if err != nil {
		return fmt.Errorf("create minio fail %s",err.Error())
	}
	s.minio = client
	return nil
}


func (s *Service) Close () error {
	return nil
}

func (s *Service) Upload(obj string, data string) error {
	isExist, err := s.minio.BucketExists(s.bucket)
	if err != nil || !isExist{
		return fmt.Errorf("Bucket %s not exist err msg: %s ",s.bucket,err.Error())
	}
	reader := strings.NewReader(data)
	_, e := s.minio.PutObject(s.bucket, obj, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/json"})
	if e != nil {
		return fmt.Errorf("Put Object %s fail ", obj)
	}
	return nil
}