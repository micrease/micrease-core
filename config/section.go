package config

import "fmt"

//常见的配置项结构定义
//服务配置项
type ServiceSection struct {
	ServiceName string `yaml:"service_name" json:"service_name"` //服务器名称
	Host        string `yaml:"host" json:"host"`                 //监听的ip,默认0.0.0.0,如果设为内存地址，则仅内网可以访问
	Port        int    `yaml:"port" json:"port"`                 //服务监听的端口号
	Env         string `yaml:"env" json:"env"`                   //服务环境,dev/test/prod
	DebugEnable bool   `yaml:"debug_enable" json:"debug_enable"` //是否开启调试选项
	Version     string `yaml:"version" json:"version"`           //当前版本号
	ConfigType  string `yaml:"config_type" json:"config_type"`   //配置类型,默认:local,nacos/consul,auto。local使用本地配置文件,nacos表示从nacos配置中心加载
}

func (s ServiceSection) ListenHost() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

//nacos配置项,以nacos为注册中心时
type NacosSection struct {
	Addrs       string `yaml:"addrs" json:"addrs"`
	NamespaceId string `yaml:"namespace_id" json:"namespace_id"`
	DataId      string `yaml:"data_id" json:"data_id"`
	Group       string `yaml:"group" json:"group"`
}

//consul配置项,以consul为注册中心时
type ConsulSection struct {
	Addrs      string `yaml:"addrs" json:"addrs"`
	ConfigPath string `yaml:"config_path" json:"config_path"`
}

//DB配置项
type DatabaseSection struct {
	Driver         string `yaml:"driver" json:"driver"` //mysql,postgres
	DataSourceName string `yaml:"dsn" json:"dsn"`
	MaxIdleConns   int    `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns   int    `yaml:"max_open_conns" json:"max_open_conns"`
	TablePrefix    string `yaml:"table_prefix" json:"table_prefix"`
	Debug          bool   `yaml:"debug" json:"debug"`
}

//DB配置项
type MutiDatabaseSection struct {
	Names string                     `yaml:"names" json:"names"`
	Dbs   map[string]DatabaseSection `yaml:"dbs" json:"dbs"`
}

//redis配置项
type RedisSection struct {
	Network  string `yaml:"network" json:"network"`
	Addr     string `yaml:"addr" json:"addr"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	Prefix   string `yaml:"prefix" json:"prefix"`
}

//redis配置项
type LoggerSection struct {
	Path   string `yaml:"path" json:"path"`     //日志存放路径
	Stdout bool   `yaml:"stdout" json:"stdout"` //控制台日志,启用后不输出到文件
	Level  string `yaml:"level" json:"level"`   //日志等级
}

type Config interface {
}

//定义一个常见的数据结构
type CommonConfig struct {
	Config
	Service  ServiceSection  `yaml:"service" json:"service"`
	Consul   ConsulSection   `yaml:"consul" yaml:"consul"`
	Database DatabaseSection `yaml:"database" json:"database"`
	Redis    RedisSection    `yaml:"redis" json:"redis"`
	Logger   LoggerSection   `yaml:"logger" json:"logger"`
}

//以上给出一个基本配置结构，如果您需要在此基础上扩展可以采用如下方式
//比如增加一个从库存的配置
//type CostumConfig struct {
//	CommonConfig
//	SlaveDB DataBaseSection
//}
