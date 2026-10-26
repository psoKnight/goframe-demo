package registry

import (
	"context"
	"strings"
	"time"

	etcdregistry "github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	defaultKeepaliveTTL = 10 * time.Second
	defaultDialTimeout  = 5 * time.Second
)

// Init enables Etcd registry when registry.etcd.address is configured.
// Otherwise GoFrame keeps its local file registry, which is suitable for single-host development.
func Init(ctx context.Context) error {
	address, err := g.Cfg().Get(ctx, "registry.etcd.address")
	if err != nil {
		return err
	}
	etcdAddress := strings.TrimSpace(address.String())
	if etcdAddress == "" {
		g.Log().Info(ctx, "Etcd registry is not configured, using GoFrame local file registry")
		return nil
	}

	keepaliveTTL, err := duration(ctx, "registry.etcd.keepaliveTTL", defaultKeepaliveTTL)
	if err != nil {
		return err
	}
	dialTimeout, err := duration(ctx, "registry.etcd.dialTimeout", defaultDialTimeout)
	if err != nil {
		return err
	}

	etcd := etcdregistry.New(etcdAddress, etcdregistry.Option{
		KeepaliveTTL: keepaliveTTL,
		DialTimeout:  dialTimeout,
	})
	// grpcx registers its resolver during package initialization. Registering the
	// registry here replaces both service registration and gRPC client discovery.
	grpcx.Resolver.Register(etcd)
	g.Log().Infof(ctx, "Etcd registry enabled: %s", etcdAddress)
	return nil
}

func duration(ctx context.Context, key string, fallback time.Duration) (time.Duration, error) {
	value, err := g.Cfg().Get(ctx, key)
	if err != nil {
		return 0, err
	}
	if strings.TrimSpace(value.String()) == "" {
		return fallback, nil
	}
	return time.ParseDuration(value.String())
}
