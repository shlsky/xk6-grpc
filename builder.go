package grpc

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	xk6_nacos "github.com/shlsky/xk6-nacos"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/codes"
	gresolver "google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
	"sync/atomic"
)

type NacosBuilder struct {
	nacosClient *xk6_nacos.NacosClient
	group       string
	nacosKey    string
}

var realNacosBuilder = NacosBuilder{}

func NewNacosBuilder() *NacosBuilder {
	return &realNacosBuilder
}

const (
	defaultWeight = 100

	// metadata key
	grpcPort       = "gRPC_port"
	envName        = "env_name"
	createTime     = "create_time"
	language       = "language"
	bsmServiceName = "bsm_service_name"
	languageGolang = "golang"
	cloudName      = "cloud_name"

	// ClusterName grpc address attribute key
	ClusterName         = "ClusterName"
	ProviderServiceName = "ProviderServiceName"

	// MyProjectEnvName env key
	MyProjectEnvName = "MY_PROJECT_ENV_NAME"

	// AwsClusterName constant
	AwsClusterName     = "aws"
	DefaultClusterName = "DEFAULT"
)

var NacosSub = make(map[string]atomic.Bool)

func (mb *NacosBuilder) FillBuilder(nacosKey string, Group string) {
	mb.nacosKey = nacosKey
	mb.group = Group
	mb.nacosClient = xk6_nacos.New()
	NacosSub[nacosKey] = atomic.Bool{}
}

type Resolver struct {
	c              naming_client.INamingClient
	group          string
	target         string
	cc             gresolver.ClientConn
	subscribeParam atomic.Pointer[vo.SubscribeParam]
	nacosKey       string
}

// Scheme for mydns
func (mb *NacosBuilder) Scheme() string {
	return "nacos"
}

func (b NacosBuilder) Build(target gresolver.Target, cc gresolver.ClientConn, opts gresolver.BuildOptions) (
	gresolver.Resolver, error) {

	GlobalResolver := &Resolver{
		c:        b.nacosClient.NacosMap[b.nacosKey],
		nacosKey: b.nacosKey,
		group:    b.group,
		// target.Endpoint fiat-go/test target.URL.Path /fiat-go/test
		// replace the first /
		target: strings.Replace(target.URL.Path, "/", "", 1),
		cc:     cc,
	}

	err := GlobalResolver.watch()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Resolver: failed to subscribe nacos: %s", err)
	}

	// warm up
	GlobalResolver.resolveNow()
	return GlobalResolver, nil
}

// ResolveNow is a no-op here.
// It's just a hint, Resolver can ignore this if it's not necessary.
func (r *Resolver) ResolveNow(options gresolver.ResolveNowOptions) {}
func (r *Resolver) resolveNow() {
	service, err := r.c.GetService(vo.GetServiceParam{
		ServiceName: r.target,
		GroupName:   r.group,
		Clusters:    []string{"ALL"},
	})

	if err == nil {
		r.updateAddress(service.Hosts, r.target)
	} else {
	}
}
func (r *Resolver) watch() error {
	if r.subscribeParam.Load() != nil {
		return nil
	}

	subscribeParam := r.buildSubscribeParam()
	v, _ := NacosSub[r.nacosKey]
	if v.CompareAndSwap(false, true) {
		return r.c.Subscribe(subscribeParam)
	}

	return nil
}

func (r *Resolver) Close() {
	if r.subscribeParam.Load() == nil {
		return
	}
	if err := r.c.Unsubscribe(r.subscribeParam.Load()); err != nil {
	}
	r.subscribeParam.Store(nil)
}
func (r *Resolver) buildSubscribeParam() *vo.SubscribeParam {
	return &vo.SubscribeParam{
		ServiceName: r.target,
		GroupName:   r.group,
		Clusters:    []string{"ALL"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			if err == nil {
				r.updateAddress(MapTo[model.SubscribeService, model.Instance](func(t model.SubscribeService) model.Instance {

					return model.Instance{
						Metadata:    t.Metadata,
						ClusterName: t.ClusterName,
						Ip:          t.Ip,
						Port:        t.Port,
						Enable:      t.Enable,
						Healthy:     t.Healthy,
					}
				}, &services), r.target)
			} else {

			}
		},
	}
}
func (r *Resolver) updateAddress(services []model.Instance, target string) {

	addrs := convertToGRPCAddresses(services, target)
	if len(addrs) > 0 {
		_ = r.cc.UpdateState(gresolver.State{Addresses: addrs})
	} else {
		// use the same error message as subscribe empty hosts error.
	}

}
func convertToGRPCAddresses(ups []model.Instance, target string) []gresolver.Address {
	var addrs []gresolver.Address

	for _, up := range ups {
		if !up.Enable || !up.Healthy {
			continue
		}

		addrs = append(addrs, convertToGRPCAddress(up.Metadata, up.ClusterName, up.Port, up.Ip, target))
	}

	return addrs
}

func convertToGRPCAddress(metadata map[string]string,
	clusterName string, port1 uint64, ip string,
	target string) gresolver.Address {
	addr := gresolver.Address{}
	var port uint64

	if len(metadata) > 0 {
		var attrs *attributes.Attributes
		for key, val := range metadata {
			// compatible with java, if metadata has gRPC_port, use it as grpc port.
			if key == grpcPort {
				if p, err := strconv.Atoi(val); err == nil {
					port = uint64(p)
				}
			}
			attrs = attrs.WithValue(key, val)
		}
		addr.Attributes = attrs
	}

	if len(clusterName) > 0 {
		attrs := addr.Attributes.WithValue(ClusterName, clusterName)
		addr.Attributes = attrs
	}

	if len(target) > 0 {
		attrs := addr.Attributes.WithValue(ProviderServiceName, target)
		addr.Attributes = attrs
	}

	if port == 0 {
		port = port1
	}

	addr.Addr = fmt.Sprintf("%s:%d", ip, port)

	return addr
}
func MapTo[T any, R any](fn func(t T) R, t *[]T) []R {
	var ans = make([]R, 0)

	for _, v := range *t {
		ans = append(ans, fn(v))
	}

	return ans
}
