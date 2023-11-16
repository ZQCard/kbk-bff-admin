module github.com/ZQCard/kbk-bff-admin

go 1.20

require (
	github.com/envoyproxy/protoc-gen-validate v0.10.1
	github.com/go-kratos/kratos/contrib/log/logrus/v2 v2.0.0-20230424154814-520b321fe99b
	github.com/go-kratos/kratos/v2 v2.6.2
	github.com/go-kratos/swagger-api v1.0.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.7.0
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.6.9
	github.com/google/wire v0.5.0
	github.com/sirupsen/logrus v1.9.0
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.uber.org/automaxprocs v1.5.2
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-kratos/grpc-gateway/v2 v2.5.1-0.20210811062259-c92d36e434b1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/glog v1.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jinzhu/copier v0.3.5
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/rakyll/statik v0.1.7 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect

)

require (
	github.com/ZQCard/kbk-administrator v1.0.0
	github.com/ZQCard/kbk-authorization v1.0.0
	github.com/ZQCard/kbk-log v1.0.0
	github.com/ZQCard/kbk-oss v1.0.0
)

replace (
	// github.com/ZQCard/kbk-administrator => ../administrator
	// github.com/ZQCard/kbk-authorization => ../authorization
	// github.com/ZQCard/kbk-log => ../kbk-log
	// github.com/ZQCard/kbk-oss => ../oss

	github.com/ZQCard/kbk-administrator => github.com/ZQCard/kbk-administrator v0.1.2
	github.com/ZQCard/kbk-authorization => github.com/ZQCard/kbk-authorization v0.1.2
	github.com/ZQCard/kbk-log => github.com/ZQCard/kbk-log v0.1.2
	github.com/ZQCard/kbk-oss => github.com/ZQCard/kbk-oss v0.1.3
)
