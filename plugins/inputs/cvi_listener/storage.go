package cvi_listener

import (
	"github.com/minio/minio-go"
)

type Storage struct {
	MinioAccess string
	MinioSecret string
	MinioURL 	string
	MinioBacket string
	client		*minio.Client
}

func (s *Storage) StartService() (error) {
	c, err := minio.New(s.MinioURL, s.MinioAccess, s.MinioSecret, false)
	if err == nil {
		s.client = c
	}

	return err
}

func (s *Storage) Upload(objectname string, filepath string) {
	s.client.FPutObject(s.MinioBacket, objectname, filepath, minio.PutObjectOptions{ContentType: "application/json"})
}