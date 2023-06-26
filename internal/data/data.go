package data

import (
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	etcdclient "go.etcd.io/etcd/client/v3"

	administratorV1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	aiV1 "github.com/ZQCard/kbk-ai/api/ai/v1"
	authorizationV1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	fileV1 "github.com/ZQCard/kbk-file/api/file/v1"
	logV1 "github.com/ZQCard/kbk-log/api/log/v1"
)

var auth *conf.Auth

func GetAuthApiKey() string {
	return auth.ApiKey
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewRegistrar,
	NewDiscovery,
	NewAdministratorRepo,
	NewAdministratorServiceClient,
	NewAuthorizationRepo,
	NewAuthorizationServiceClient,
	NewApiLogRepo,
	NewApiLogClient,
	NewFileRepo,
	NewFileClient,
	NewRedisClient,
	NewAIRepo,
	NewAIClient,
)

// Data .
type Data struct {
	cfg                 *conf.Bootstrap
	logger              *log.Helper
	rdb                 *redis.Client
	administratorClient administratorV1.AdministratorServiceClient
	authorizationClient authorizationV1.AuthorizationServiceClient
	logClient           logV1.LogServiceClient
	fileClient          fileV1.FileServiceClient
	aiClient            aiV1.AIServiceClient
}

func NewData(
	cfg *conf.Bootstrap,
	redisCli *redis.Client,
	logger log.Logger,
	administratorClient administratorV1.AdministratorServiceClient,
	authorizationClient authorizationV1.AuthorizationServiceClient,
	logClient logV1.LogServiceClient,
	fileClient fileV1.FileServiceClient,
	aiClient aiV1.AIServiceClient,
) (*Data, func(), error) {
	auth = cfg.Auth
	logs := log.NewHelper(log.With(logger, "module", "kratos-base-layout/data"))
	cleanup := func() {
		logs.Info("closing the data resources")
	}

	return &Data{
		logger:              logs,
		cfg:                 cfg,
		administratorClient: administratorClient,
		authorizationClient: authorizationClient,
		logClient:           logClient,
		rdb:                 redisCli,
		fileClient:          fileClient,
		aiClient:            aiClient,
	}, cleanup, nil
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	point := conf.Etcd.Address
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{point},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)
	return r
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	point := conf.Etcd.Address
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{point},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)
	return r
}

func NewRedisClient(conf *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	err := client.Ping().Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}
