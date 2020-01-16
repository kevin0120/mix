package utils

// Service represents a service attached to the server.
type ICommonService interface {
	Open() error
	Close() error
}
