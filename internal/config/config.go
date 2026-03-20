package config

import "time"

type Config struct {
	GRPC     GRPCConfig     `yaml:"grpc"`
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
	MinIO    MinIOConfig    `yaml:"minio"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
}

type HTTPConfig struct {
	Port int `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"PG_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PG_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"PG_USER" env-default:"grusha"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-default:"grusha_secret"`
	DBName   string `yaml:"db_name" env:"PG_DBNAME" env-default:"grusha"`
	SSLMode  string `yaml:"ssl_mode" env:"PG_SSLMODE" env-default:"disable"`
}

func (c PostgresConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password +
		"@" + c.Host + ":" + itoa(c.Port) +
		"/" + c.DBName + "?sslmode=" + c.SSLMode
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	DB       int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
}

type MinIOConfig struct {
	Endpoint  string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	AccessKey string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	Bucket    string `yaml:"bucket" env:"MINIO_BUCKET" env-default:"grusha"`
	UseSSL    bool   `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
}

type JWTConfig struct {
	Secret          string        `yaml:"secret" env:"JWT_SECRET" env-default:"super-secret-key-change-me"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"JWT_REFRESH_TTL" env-default:"720h"`
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
