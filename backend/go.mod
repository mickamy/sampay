module mickamy.com/sampay

go 1.23.4

require (
	buf.build/gen/go/mickamy/sampay/connectrpc/go v1.18.1-20250112000630-58e1ac5e43ed.1
	buf.build/gen/go/mickamy/sampay/protocolbuffers/go v1.36.2-20250112000630-58e1ac5e43ed.1
	connectrpc.com/connect v1.18.1
	github.com/DATA-DOG/go-txdb v0.2.0
	github.com/alicebob/miniredis/v2 v2.34.0
	github.com/aws/aws-sdk-go-v2 v1.32.8
	github.com/aws/aws-sdk-go-v2/config v1.28.10
	github.com/aws/aws-sdk-go-v2/service/ssm v1.56.4
	github.com/brianvoe/gofakeit/v7 v7.1.2
	github.com/caarlos0/env/v11 v11.3.1
	github.com/docker/docker v27.4.1+incompatible
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang-migrate/migrate/v4 v4.18.1
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/mickamy/slogger v0.0.6
	github.com/nicksnyder/go-i18n/v2 v2.4.1
	github.com/oklog/ulid/v2 v2.1.0
	github.com/redis/go-redis/v9 v9.7.0
	github.com/rs/cors v1.11.1
	github.com/stretchr/testify v1.10.0
	github.com/testcontainers/testcontainers-go v0.35.0
	go.uber.org/mock v0.5.0
	golang.org/x/crypto v0.32.0
	golang.org/x/net v0.34.0
	golang.org/x/text v0.21.0
	google.golang.org/protobuf v1.36.2
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.5.11
	gorm.io/gorm v1.25.12
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.2-20241127180247-a33202765966.1 // indirect
	dario.cat/mergo v1.0.1 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20240806141605-e8a1dd7889d6 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.51 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.6 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.9.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
