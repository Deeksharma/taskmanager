package server

type Server interface {
	Start()
	Stop() <-chan struct{}
}
