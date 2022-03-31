package config

import (
	"fmt"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"io/ioutil"
	"strconv"
	"strings"
)

//当nacos做为配置中心时，配置需要从nacos加载
var nacosClientConfig config_client.IConfigClient

func LoadConfigContent(nacos NacosSection) (string, error) {
	//从nacos配置中心拉取业务配置
	content, err := loadConfigFromNacos(nacos)
	if err == nil {
		//监听nacos配置中心数据变化
		listenChange(nacos)
	}
	return content, err
}

func loadConfigFromNacos(nacos NacosSection) (string, error) {
	nacosClientConfig := GetConfigClient(nacos)
	dataId := nacos.DataId
	group := nacos.Group
	content, err := nacosClientConfig.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group})

	if err != nil {
		log.Errorf("Nacos配置中心连接失败,{env}-nacos.ini配置不正确,"+err.Error()+"!!,NacosAddrs=%s,NamespaceId=%s,DataId=%s,Group=%s",
			nacos.Addrs,
			nacos.NamespaceId,
			nacos.DataId,
			nacos.Group,
		)
		//os.Exit(0)
		return "", err
	}

	log.Debug("\n", content)
	return content, err
}

func writeTo(filename string, text string) error {
	var content = []byte(text)
	err := ioutil.WriteFile(filename, content, 0666) //写入文件(字节数组)
	return err
}

//on_change
func listenChange(nacos NacosSection) {
	dataId := nacos.DataId
	group := nacos.Group
	err := nacosClientConfig.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Info("ListenChange Run")
			fmt.Println("ListenChange Run", data)
			//LoadConfigContent(nacos)
		},
	})
	log.Error(err)
}

func NewClientConfig(nacos NacosSection) constant.ClientConfig {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(nacos.NamespaceId), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		//constant.WithLogDir("./nacos/log"),
		//constant.WithCacheDir("./nacos/cache"),
		//constant.WithRotateTime("1h"),
		//constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)
	return clientConfig
}

func NewServerConfig(nacos NacosSection) []constant.ServerConfig {
	addrs := nacos.Addrs
	arr := strings.Split(addrs, ":")

	if arr[0] == "" {
		log.Error("nacos addr is empty")
	}
	addr := arr[0]
	port := 8848
	if arr[1] != "" {
		intNum, _ := strconv.Atoi(arr[1])
		port = int(uint64(intNum))
	}

	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			addr,
			uint64(port),
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	return serverConfigs
}

func GetNamingConfigClient(nacos NacosSection) naming_client.INamingClient {
	clientConfig := NewClientConfig(nacos)
	serverConfigs := NewServerConfig(nacos)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Error("nacos error", err)
	}
	return namingClient
}

func GetConfigClient(nacos NacosSection) config_client.IConfigClient {
	clientConfig := NewClientConfig(nacos)
	serverConfigs := NewServerConfig(nacos)
	var err error
	nacosClientConfig, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Error("nacos error", err)
	}
	return nacosClientConfig
}
