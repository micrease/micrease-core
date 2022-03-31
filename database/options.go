package database

import "time"

//数据库连接信息
type Options struct {
	Name         string
	Driver       string
	Dsn          string
	OrmType      string
	Timeout      time.Duration
	MaxIdleConns int32
	MaxOpenConns int32
}

type Option func(*Options)

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Driver(driver string) Option {
	return func(o *Options) {
		o.Driver = driver
	}
}

func Dsn(dsn string) Option {
	return func(o *Options) {
		o.Dsn = dsn
	}
}

func Timeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}
