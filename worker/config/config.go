package config

type Config struct {
	DBConn    string
	RedisAddr string
}

func Get() *Config {
	return &Config{
		DBConn:    "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable",
		RedisAddr: "localhost:6379",
	}
}
