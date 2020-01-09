package minio

import "github.com/masami10/rush/services/storage"

type IStorageService interface {
	ListUnuploadCurves() ([]storage.Curves, error)
	UpdateCurve(curve *storage.Curves) (*storage.Curves, error)
}
