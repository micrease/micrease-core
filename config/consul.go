package config

import (
	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	cconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
)

func LoadConsulConfig(sec ConsulSection, dest interface{}) error {
	//添加配置中心
	consulClient, err := api.NewClient(&api.Config{
		Address: sec.Addrs,
	})
	if err != nil {
		panic(err)
	}
	cs, err := consul.New(consulClient, consul.WithPath(sec.ConfigPath))
	//consul中需要标注文件后缀，kratos读取配置需要适配文件后缀
	//The file suffix needs to be marked, and kratos needs to adapt the file suffix to read the configuration.
	if err != nil {
		panic(err)
	}
	c := cconfig.New(cconfig.WithSource(cs))
	err = c.Load()
	if err != nil {
		return err
	}
	return c.Scan(dest)
}
