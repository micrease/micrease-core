package config

//常见的配置项结构定义
//服务配置项
type ServiceSection struct {
	ServiceName string `yaml:"service_name"` //服务器名称
	Host        string `yaml:"host"`         //监听的ip,默认0.0.0.0,如果设为内存地址，则仅内网可以访问
	Port        string `yaml:"port"`         //服务监听的端口号
	Env         string `yaml:"env"`          //服务环境,dev/test/prod
	DebugEnable bool   `yaml:"debug_enable"` //是否开启调试选项
	Version     string `yaml:"version"`      //当前版本号
	ConfigType  string `yaml:"config_type"`  //配置类型,默认:local,nacos。local使用本地配置文件,nacos表示从nacos配置中心加载
}

//nacos配置项
type NacosSection struct {
	Addrs       string `yaml:"addrs"`
	NamespaceId string `yaml:"namespace_id"`
	DataId      string `yaml:"data_id"`
	Group       string `yaml:"group"`
}

//DB配置项
type DatabaseSection struct {
	Driver         string `yaml:"driver"` //mysql,postgres
	DataSourceName string `yaml:"dsn"`
	MaxIdleConns   int32  `yaml:"max_idle_conns"`
	MaxOpenConns   int32  `yaml:"max_open_conns"`
	TablePrefix    string `yaml:"table_prefix"`
	Debug          bool   `yaml:"debug"`
}

//DB配置项
type MutiDatabaseSection struct {
	Names string                     `yaml:"names"`
	Dbs   map[string]DatabaseSection `yaml:"dbs"`
}

//redis配置项
type RedisSection struct {
	NetWork  string `yaml:"net_work"`
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

//redis配置项
type LoggerSection struct {
	Path   string `yaml:"path"`   //日志存放路径
	Stdout bool   `yaml:"stdout"` //控制台日志,启用后不输出到文件
	Level  string `yaml:"level"`  //日志等级
}

//定义一个常见的数据结构
type CommonConfig struct {
	Service  ServiceSection  `yaml:"service"`
	Nacos    NacosSection    `yaml:"nacos"`
	Database DatabaseSection `yaml:"database"`
	Redis    RedisSection    `yaml:"redis"`
	Logger   LoggerSection   `yaml:"logger"`
}

//以上给出一个基本配置结构，如果您需要在此基础上扩展可以采用如下方式
//比如增加一个从库存的配置
//type CostumConfig struct {
//	CommonConfig
//	SlaveDB DataBaseSection
//}
