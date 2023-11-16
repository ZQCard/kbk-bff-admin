package data

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	"github.com/ZQCard/kbk-bff-admin/internal/conf"

	administratorV1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	authorizationV1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	logV1 "github.com/ZQCard/kbk-log/api/log/v1"
	ossV1 "github.com/ZQCard/kbk-oss/api/oss/v1"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewAdministratorRepo,
	NewAdministratorServiceClient,
	NewAuthorizationRepo,
	NewAuthorizationServiceClient,
	NewApiLogRepo,
	NewApiLogClient,
	NewOssRepo,
	NewOssClient,
	NewRedisClient,
)

// Data .
type Data struct {
	cfg                 *conf.Bootstrap
	logger              *log.Helper
	rds                 *redis.Client
	administratorClient administratorV1.AdministratorServiceClient
	authorizationClient authorizationV1.AuthorizationServiceClient
	logClient           logV1.LogServiceClient
	ossClient           ossV1.OSSClient
}

func NewData(
	cfg *conf.Bootstrap,
	redisCli *redis.Client,
	logger log.Logger,
	administratorClient administratorV1.AdministratorServiceClient,
	authorizationClient authorizationV1.AuthorizationServiceClient,
	logClient logV1.LogServiceClient,
	ossClient ossV1.OSSClient,
) (*Data, func(), error) {
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
		rds:                 redisCli,
		ossClient:           ossClient,
	}, cleanup, nil
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
