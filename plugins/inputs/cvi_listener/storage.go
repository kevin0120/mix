package cvi_listener

import (
	"github.com/minio/minio-go"
	"strings"
)

type Storage struct {
	MinioAccess string
	MinioSecret string
	MinioURL 	string
	MinioBacket string
	client		*minio.Client
}

//func (s *Storage) StartService() (error) {
//
//}

func (s *Storage) Upload(objectname string, data string) (error) {
	c, err := minio.New(s.MinioURL, s.MinioAccess, s.MinioSecret, false)
	if err == nil {
		s.client = c
	} else {
		return err
	}

	reader := strings.NewReader(data)
	_, e := s.client.PutObject(s.MinioBacket, objectname, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/json"})
	if e != nil {
		return e
	}

	return nil
}