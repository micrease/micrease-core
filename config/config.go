package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	_map "micrease-core/map"
	"sync"
)

//定义一个全局的key-value结构的配置表,用来装载或加载配置项
//比如初始化服务时将配置信息加载进该配置表。在应用的其它任何位置可以通过key获取配置
type ConfigMap struct {
	_map.Map
}

var once sync.Once
var gConfig *ConfigMap

func GetDefaultConfig() *ConfigMap {
	once.Do(func() {
		gConfig = newConfig()
	})
	return gConfig
}

func newConfig() *ConfigMap {
	gConfig = new(ConfigMap)
	m := make(map[string]interface{})
	gConfig.NewMap(m)
	return gConfig
}

//config.Set("xx","yyy")
func Set(key string, value interface{}) {
	GetDefaultConfig().Set(key, value)
}

//config.Get("xx")
func Get(key string) {
	GetDefaultConfig().Get(key)
}
func ConvertToMap(cfg *ini.File) {
	for _, section := range cfg.Sections() {
		for _, key := range section.Keys() {

			k := section.Name() + "." + key.Name()
			if section.Name() == ini.DefaultSection {
				k = key.Name()
			}
			fmt.Println("config:" + k + "=>" + key.Value())
			Set(k, key.Value())
		}
	}
}
