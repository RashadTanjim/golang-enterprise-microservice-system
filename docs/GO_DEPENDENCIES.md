# Go Dependencies Reference

This document describes every module listed in `go.mod` and why it exists in this codebase. Direct dependencies are used explicitly in code. Indirect dependencies are pulled in transitively by direct dependencies (frameworks, tooling, database drivers, etc.).

## Direct Dependencies

| Module | Purpose in this project |
| --- | --- |
| github.com/gin-gonic/gin | HTTP framework used by all services. |
| github.com/golang-jwt/jwt/v5 | JWT signing and verification for auth. |
| github.com/golang-migrate/migrate/v4 | Database migrations executed by the migration service. |
| github.com/google/uuid | Request ID generation middleware. |
| github.com/joho/godotenv | Loads `.env` files in local/dev. |
| github.com/prometheus/client_golang | Prometheus metrics instrumentation. |
| github.com/redis/go-redis/v9 | Redis client for caching in services. |
| github.com/sony/gobreaker | Circuit breaker in order service. |
| github.com/stretchr/testify | Testing assertions and mocks. |
| github.com/swaggo/files | Swagger UI static assets. |
| github.com/swaggo/gin-swagger | Swagger UI handler for Gin. |
| github.com/swaggo/swag | Swagger doc generation support. |
| go.uber.org/zap | Structured JSON logging. |
| golang.org/x/time | Token bucket rate limiting. |
| gorm.io/driver/postgres | PostgreSQL driver for GORM. |
| gorm.io/driver/sqlite | SQLite driver for unit tests. |
| gorm.io/gorm | ORM used by repositories. |

## Indirect Dependencies

| Module | Why it is present |
| --- | --- |
| github.com/KyleBanks/depth | Swagger generator dependency graph utility (swag). |
| github.com/PuerkitoBio/purell | URL normalization for Swagger tooling. |
| github.com/PuerkitoBio/urlesc | URL escaping used by Swagger tooling. |
| github.com/beorn7/perks | Prometheus client metrics helpers. |
| github.com/bytedance/sonic | JSON acceleration used by Gin. |
| github.com/bytedance/sonic/loader | Sonic runtime loader. |
| github.com/cespare/xxhash/v2 | Fast hashing used by Prometheus client. |
| github.com/cloudwego/base64x | Base64 helpers used by Sonic. |
| github.com/cloudwego/iasm | Assembly optimizations used by Sonic. |
| github.com/davecgh/go-spew | Deep pretty-printer used by testify. |
| github.com/dgryski/go-rendezvous | Rendezvous hashing used by go-redis. |
| github.com/gabriel-vasile/mimetype | Content-type detection in Gin. |
| github.com/gin-contrib/sse | Server-sent events utilities in Gin. |
| github.com/go-openapi/jsonpointer | JSON pointer handling for Swagger tooling. |
| github.com/go-openapi/jsonreference | JSON reference handling for Swagger tooling. |
| github.com/go-openapi/spec | OpenAPI spec structures used by swag. |
| github.com/go-openapi/swag | OpenAPI helpers used by swag. |
| github.com/go-playground/locales | Locale data for validator. |
| github.com/go-playground/universal-translator | i18n support for validator. |
| github.com/go-playground/validator/v10 | Struct validation used by Gin binding. |
| github.com/goccy/go-json | JSON encoder/decoder used by Gin. |
| github.com/jackc/pgpassfile | Postgres password file support in pgx. |
| github.com/jackc/pgservicefile | Postgres service file support in pgx. |
| github.com/jackc/pgx/v5 | Postgres driver used by GORM postgres. |
| github.com/jackc/puddle/v2 | Connection pool used by pgx. |
| github.com/jinzhu/inflection | GORM naming/inflection utilities. |
| github.com/jinzhu/now | GORM time helpers. |
| github.com/josharian/intern | String interning for Swagger tooling. |
| github.com/json-iterator/go | JSON library used by Gin. |
| github.com/klauspost/cpuid/v2 | CPU feature detection used by Sonic. |
| github.com/leodido/go-urn | URN parsing used by validator. |
| github.com/lib/pq | Postgres driver used by migrate. |
| github.com/mailru/easyjson | JSON helpers used by Swagger tooling. |
| github.com/mattn/go-isatty | TTY detection used by logging deps. |
| github.com/mattn/go-sqlite3 | SQLite driver for tests. |
| github.com/modern-go/concurrent | JSON iterator dependency. |
| github.com/modern-go/reflect2 | JSON iterator dependency. |
| github.com/pelletier/go-toml/v2 | TOML parsing used by tooling. |
| github.com/pmezard/go-difflib | Diff utilities used by testify. |
| github.com/prometheus/client_model | Prometheus data model types. |
| github.com/prometheus/common | Prometheus shared helpers. |
| github.com/prometheus/procfs | Process metrics for Prometheus client. |
| github.com/stretchr/objx | Testify helper library. |
| github.com/twitchyliquid64/golang-asm | Assembly helpers used by Sonic. |
| github.com/ugorji/go/codec | JSON/codec used by Gin. |
| go.uber.org/multierr | Error aggregation used by Zap. |
| golang.org/x/arch | Low-level architecture helpers. |
| golang.org/x/crypto | Crypto utilities used by deps. |
| golang.org/x/net | Network helpers used by deps. |
| golang.org/x/sync | Sync primitives used by deps. |
| golang.org/x/sys | OS syscalls used by deps. |
| golang.org/x/text | Text processing for validators and URL tooling. |
| golang.org/x/tools | Tooling libs used by swag. |
| google.golang.org/protobuf | Protobuf types used by Prometheus. |
| gopkg.in/yaml.v2 | YAML parsing used by tooling. |
| gopkg.in/yaml.v3 | YAML parsing used by tooling. |
