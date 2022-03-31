package bootstrap

import (
	"sync"
)

const (
	ServiceTypeHttp = "web"
	ServiceTypeRpc = "rpc"
)

type MicroService interface {
	Name() string

	Init(...Option)
	// Options returns the current options
	Options() Options
	// Run the service
	Start() error
	// The service implementation
	String() string
}

//继承MicroService
type microService struct {
	opts Options
	once sync.Once
}

func (m microService) Init(...Option) {
}

func (m microService) Name() string {
	return ""
}

func (m microService) Options() Options {
	return newOptions()
}

func (m microService) Start() error {
	return nil
}

func (m microService) String() string {
	return ""
}

func NewMicroService(opts ...Option) MicroService {
	return newMicroBoot(opts...)
}

func newMicroBoot(opts ...Option) MicroService {
	service := new(microService)
	options := newOptions(opts...)
	// set opts
	service.opts = options
	return service
}
