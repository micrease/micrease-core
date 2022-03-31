package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//配置加载器
//加载逻辑为:
//1,加载本地./resources配置项
//2,尝试加载配置中心配置项,如果配置中心可用，则使用配置中心的配置

var (
	ResourcesPath = flag.String("resources", "./resources", "-config resources path")
	Env           = flag.String("env", "dev", "-env:dev/test/uat/prod")
)

const (
	//从配置中心加载配置,本地配置中至少需要nacos配置块
	ConfigTypeNacos = "nacos"
	//从配置文件加载
	ConfigTypeLocal = "local"
	//自动加载配置策略
	ConfigTypeAuto = "auto"
)

var commonConfig = new(CommonConfig)

func GetCommonConfig() *CommonConfig {
	return commonConfig
}

//1,加载配置到一个常规的配置结构体中
func LoadConfig() *CommonConfig {
	LoadConfigTo(commonConfig)
	return commonConfig
}

//2,加载配置到自定义的配置结构体中
func LoadConfigTo(dest interface{}) error {
	//首先从本地加载commonConfig,获取nacos信息
	err := LoadLocalConfigTo(commonConfig)
	if err != nil {
		return err
	}

	if commonConfig.Service.ConfigType == ConfigTypeNacos && len(commonConfig.Nacos.Addrs) > 0 {
		//从配置中心加载配置
		content, err := LoadConfigContent(commonConfig.Nacos)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal([]byte(content), dest); err != nil {
			return err
		}
		return nil
	}

	err = LoadLocalConfigTo(dest)
	if err != nil {
		return err
	}
	return nil
}

func LoadLocalConfigTo(dest interface{}) error {
	path := getLocalConfigFile()
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(content, dest); err != nil {
		return err
	}
	return nil
}

//获取当地配置文件地址
func getLocalConfigFile() string {
	resourcesPath := *ResourcesPath
	env := *Env
	path := resourcesPath + "/config-" + env + ".yml"
	return path
}

//这是一个缓存配置
func getLocalTempFile() string {
	resourcesPath := *ResourcesPath
	path := resourcesPath + "/cache.yml"
	return path
}
