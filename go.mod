module github.com/iRankHub/backend

go 1.23

require (
	github.com/aws/aws-sdk-go-v2 v1.32.2
	github.com/aws/aws-sdk-go-v2/config v1.27.43
	github.com/aws/aws-sdk-go-v2/credentials v1.17.41
	github.com/aws/aws-sdk-go-v2/service/s3 v1.65.3
	github.com/go-webauthn/webauthn v0.10.2
	github.com/golang-migrate/migrate/v4 v4.17.1
	github.com/jackc/pgx/v5 v5.5.5
	github.com/lib/pq v1.10.9
	github.com/o1egl/paseto v1.0.0
	github.com/pquerna/otp v1.4.0
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sqlc-dev/pqtype v0.3.0
	golang.org/x/crypto v0.31.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	k8s.io/api v0.31.1
	k8s.io/apimachinery v0.31.1
	k8s.io/client-go v0.31.1
	k8s.io/metrics v0.31.1
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/chacha20poly1305 v0.0.0-20170617001512-233f39982aeb // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.2 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/docker v27.1.1+incompatible // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-webauthn/x v0.1.9 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.52.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
