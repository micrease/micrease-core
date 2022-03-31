package bootstrap

import (
	"github.com/micro/go-micro/v2"
	"micro-core/config"
	"micro-core/database"
)

const (
	SERVICE_RPC  = "rpc"
	SERVICE_HTTP = "http"
)

type Options struct {
	Name        string
	ServiceType string
	Orm         database.Orm
	Service     micro.Service
	Config      config.Config
}
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:        "",
		ServiceType: SERVICE_HTTP,
		Orm:         database.DefaultOrm,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func ServiceType(st string) Option {
	return func(o *Options) {
		o.ServiceType = st
	}
}

func Orm(orm database.Orm) Option {
	return func(o *Options) {
		o.Orm = orm
	}
}

func Service(service micro.Service) Option {
	return func(o *Options) {
		o.Service = service
	}
}

func Config(config config.Config) Option {
	return func(o *Options) {
		o.Config = config
	}
}
