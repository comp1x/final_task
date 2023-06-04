package config

type Config struct {
	Customer
	DB
	Restaurant
}

type Customer struct {
	HTTPAddr string `env:"GATEWAY_HTTP_ADDR" envDefault:":13101"`
	GRPCAddr string `env:"GATEWAY_GRPC_ADDR" envDefault:":13102"`
}

type DB struct {
	PgPort   string `env:"PG_PORT" envDefault:"5454"`
	PgHost   string `env:"PG_HOST" envDefault:"0.0.0.0"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"db"`
	PgUser   string `env:"PG_USER" envDefault:"db"`
	PgPwd    string `env:"PG_PWD" envDefault:"db"`
}

type Restaurant struct {
	GRPCAddr string `env:"GATEWAY_GRPC_ADDR_RESTAURANT" envDefault:":13104"`
}
