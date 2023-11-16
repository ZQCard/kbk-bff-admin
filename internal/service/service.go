package service

import (
	v1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminInterface)

type AdminInterface struct {
	v1.UnimplementedAdminServer
	log               *log.Helper
	administratorRepo *data.AdministratorRepo
	authorizationRepo *data.AuthorizationRepo
	apiLogRepo        *data.ApiLogRepo
	ossRepo           *data.OssRepo
}

func NewAdminInterface(
	logger log.Logger,
	administratorRepo *data.AdministratorRepo,
	authorizationRepo *data.AuthorizationRepo,
	apiLogRepo *data.ApiLogRepo,
	ossRepo *data.OssRepo,
) *AdminInterface {
	return &AdminInterface{
		log:               log.NewHelper(log.With(logger, "module", "bff-admin/interface")),
		administratorRepo: administratorRepo,
		apiLogRepo:        apiLogRepo,
		ossRepo:           ossRepo,
		authorizationRepo: authorizationRepo,
	}
}
