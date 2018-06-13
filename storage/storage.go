package rush_storage

import (
	"github.com/minio/minio-go"
	"strings"
)

type StorageConf struct {
	URL    string `yaml:"url"`
	Backet string `yaml:"backet"`
	Access string `yaml:"access"`
	Secret string `yaml:"secret"`
}

type Storage struct {
	Conf   *StorageConf
	client *minio.Client
}

func (s *Storage) Upload(objectname string, data string) error {
	c, err := minio.New(s.Conf.URL, s.Conf.Access, s.Conf.Secret, false)
	if err == nil {
		s.client = c
	} else {
		return err
	}

	reader := strings.NewReader(data)
	_, e := s.client.PutObject(s.Conf.Backet, objectname, reader, reader.Size(), minio.PutObjectOptions{ContentType: "application/json"})
	if e != nil {
		return e
	}

	return nil
}
