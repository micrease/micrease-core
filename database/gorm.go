package database

import (
	"time"
)

var (
	DefaultDriver = "mysql"
	DefaultORM    = "gorm"
)

type Gorm struct {
	opts Options
}

func (d *Gorm) Init(...Option) error {
	return nil
}

func (d *Gorm) Connect() error {
	return nil
}

func (d *Gorm) Get() interface{} {
	return nil
}

func (d *Gorm) String() string {
	return DefaultDriver
}

func newOrm(opts ...Option) Orm {
	options := Options{
		Driver:       DefaultDriver,
		OrmType:      DefaultORM,
		Timeout:      time.Millisecond * 100,
		MaxIdleConns: 10,
	}

	for _, o := range opts {
		o(&options)
	}

	return &Gorm{
		opts: options,
	}
}

func NewOrm(opts ...Option) Orm {
	return newOrm(opts...)
}
