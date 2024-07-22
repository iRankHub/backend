# Dependencies

This document lists the dependencies used in the iRankHub backend application.

## Go Packages

The following Go packages are used in the project:

- `aidanwoods.dev/go-paseto`: A library for generating and verifying PASETO tokens.
- `github.com/fsnotify/fsnotify`: A library for file system notifications.
- `github.com/golang/mock`: A mocking framework for Go.
- `github.com/hashicorp/hcl`: A library for parsing HCL (HashiCorp Configuration Language) files.
- `github.com/jackc/pgpassfile`: A library for handling PostgreSQL password files.
- `github.com/jackc/pgservicefile`: A library for handling PostgreSQL service files.
- `github.com/jackc/pgx/v5`: A PostgreSQL driver and toolkit for Go.
- `github.com/magiconair/properties`: A library for reading and writing properties files.
- `github.com/mitchellh/mapstructure`: A library for decoding generic map values into structs.
- `github.com/pelletier/go-toml/v2`: A library for parsing TOML files.
- `github.com/rs/cors`: A library for handling Cross-Origin Resource Sharing (CORS).
- `github.com/sagikazarmark/locafero`: A library for file system abstraction.
- `github.com/sagikazarmark/slog-shim`: A library for structured logging.
- `github.com/sourcegraph/conc`: A library for concurrent programming utilities.
- `github.com/spf13/afero`: A library for file system abstraction.
- `github.com/spf13/cast`: A library for type casting.
- `github.com/spf13/pflag`: A library for parsing command-line flags.
- `github.com/spf13/viper`: A library for configuration management.
- `github.com/stretchr/testify`: A library for writing test assertions and mocks.
- `github.com/subosito/gotenv`: A library for loading environment variables from files.
- `go.uber.org/atomic`: A library for atomic operations.
- `go.uber.org/multierr`: A library for error handling.
- `golang.org/x/crypto`: A library for cryptographic primitives.
- `golang.org/x/exp`: A library for experimental packages.
- `golang.org/x/mod`: A library for module version management.
- `golang.org/x/net`: A library for network utilities.
- `golang.org/x/sys`: A library for system-level operations.
- `golang.org/x/text`: A library for text processing.
- `golang.org/x/tools`: A library for Go tools.
- `google.golang.org/genproto/googleapis/rpc`: A library for generated protocol buffers for Google APIs.
- `google.golang.org/grpc`: A library for gRPC communication.
- `google.golang.org/protobuf`: A library for protocol buffers.
- `gopkg.in/ini.v1`: A library for parsing INI files.
- `gopkg.in/yaml.v3`: A library for parsing YAML files.

For the specific versions used, please refer to the `go.mod` file.

## Updating Dependencies

To update the dependencies, run the following command:

```bash
go get -u ./...
```