package nacos

import (
	"context"
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go-micro.dev/v4/cmd"
	"go-micro.dev/v4/registry"
	mnet "go-micro.dev/v4/util/net"
	"net"
	"strconv"
	"time"
)

type nacosRegistry struct {
	namingClient naming_client.INamingClient
	opts         registry.Options
}

func init() {
	cmd.DefaultRegistries["nacos"] = NewRegistry
}

func DefaultListOption() registry.ListOption {
	param := vo.GetAllServiceInfoParam{
		NameSpace: "public",
		GroupName: "DEFAULT_GROUP",
	}
	ctx := context.WithValue(context.Background(), "get_all_service_info_param", param)
	return registry.ListContext(ctx)
}

func GetListOption(namespace string, groupName string) registry.ListOption {
	param := vo.GetAllServiceInfoParam{
		NameSpace: namespace,
		GroupName: groupName,
	}
	ctx := context.WithValue(context.Background(), "get_all_service_info_param", param)
	return registry.ListContext(ctx)
}

func getNodeIpPort(s *registry.Service) (host string, port int, err error) {
	if len(s.Nodes) == 0 {
		return "", 0, errors.New("you must deregister at least one node")
	}
	node := s.Nodes[0]
	host, pt, err := net.SplitHostPort(node.Address)
	if err != nil {
		return "", 0, err
	}
	port, err = strconv.Atoi(pt)
	if err != nil {
		return "", 0, err
	}
	return
}

func NewClientConfig() constant.ClientConfig {
	clientConfig := *constant.NewClientConfig(
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir("./runtime/cache"),
		constant.WithLogDir("./runtime/log"),
		constant.WithLogLevel("debug"),
	)
	return clientConfig
}

func configure(c *nacosRegistry, opts ...registry.Option) error {
	// set opts
	for _, o := range opts {
		o(&c.opts)
	}
	if c.opts.Context != nil {
		if client, ok := c.opts.Context.Value("naming_client").(naming_client.INamingClient); ok {
			c.namingClient = client
			return nil
		}
	}
	clientConfig := NewClientConfig()
	serverConfigs := make([]constant.ServerConfig, 0)
	contextPath := "/nacos"
	// iterate the options addresses
	for _, address := range c.opts.Addrs {
		// check we have a port
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        8848,
				ContextPath: contextPath,
			})
		} else if err == nil {
			p, err := strconv.ParseUint(port, 10, 64)
			if err != nil {
				continue
			}
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        p,
				ContextPath: contextPath,
			})
		}
	}

	if c.opts.Timeout == 0 {
		c.opts.Timeout = time.Second * 1
	}
	clientConfig.TimeoutMs = uint64(c.opts.Timeout.Milliseconds())
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return err
	}
	c.namingClient = client
	return nil
}

func (c *nacosRegistry) Init(opts ...registry.Option) error {
	return configure(c, opts...)
}

func (c *nacosRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {
	var options registry.DeregisterOptions
	for _, o := range opts {
		o(&options)
	}
	withContext := false
	param := vo.DeregisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("deregister_instance_param").(vo.DeregisterInstanceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		host, port, err := getNodeIpPort(s)
		if err != nil {
			return err
		}
		param.Ip = host
		param.Port = uint64(port)
		param.ServiceName = s.Name
	}

	_, err := c.namingClient.DeregisterInstance(param)
	return err
}

func (c *nacosRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	withContext := false
	param := vo.RegisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("register_instance_param").(vo.RegisterInstanceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		host, port, err := getNodeIpPort(s)
		if err != nil {
			return err
		}

		if len(s.Nodes[0].Metadata) == 0 {
			s.Nodes[0].Metadata = make(map[string]string)
		}

		s.Nodes[0].Metadata["version"] = s.Version
		param.Ip = host
		param.Port = uint64(port)
		param.Metadata = s.Nodes[0].Metadata
		param.ServiceName = s.Name
		param.Enable = true
		param.Healthy = true
		param.Weight = 1.0
		param.Ephemeral = true
	}
	_, err := c.namingClient.RegisterInstance(param)
	return err
}

func (c *nacosRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	var options registry.GetOptions
	for _, o := range opts {
		o(&options)
	}
	withContext := false
	param := vo.GetServiceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("select_instances_param").(vo.GetServiceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		param.ServiceName = name
	}
	service, err := c.namingClient.GetService(param)
	if err != nil {
		return nil, err
	}
	services := make([]*registry.Service, 0)
	for _, v := range service.Hosts {
		nodes := make([]*registry.Node, 0)
		nodes = append(nodes, &registry.Node{
			Id:       v.InstanceId,
			Address:  mnet.HostPort(v.Ip, v.Port),
			Metadata: v.Metadata,
		})
		s := registry.Service{
			Name:     v.ServiceName,
			Version:  v.Metadata["version"],
			Metadata: v.Metadata,
			Nodes:    nodes,
		}
		services = append(services, &s)
	}

	return services, nil
}

func (c *nacosRegistry) ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	var options registry.ListOptions
	for _, o := range opts {
		o(&options)
	}

	param := vo.GetAllServiceInfoParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("get_all_service_info_param").(vo.GetAllServiceInfoParam); ok {
			param = p
		}
	}

	if param.PageNo == 0 {
		param.PageNo = 1
	}

	if param.PageSize == 0 {
		param.PageSize = 1024
	}

	services, err := c.namingClient.GetAllServicesInfo(param)
	if err != nil {
		return nil, err
	}
	var registryServices []*registry.Service
	for _, v := range services.Doms {
		registryServices = append(registryServices, &registry.Service{Name: v})
	}
	return registryServices, nil
}

func (c *nacosRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return NewNacosWatcher(c, opts...)
}

func (c *nacosRegistry) String() string {
	return "nacos"
}

func (c *nacosRegistry) Options() registry.Options {
	return c.opts
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	nacos := &nacosRegistry{
		opts: registry.Options{
			Context: context.Background(),
		},
	}
	configure(nacos, opts...)
	return nacos
}
