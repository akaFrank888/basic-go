package config

var Config = config{
	DB: DBConfig{
		DSN: "root:123456@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
