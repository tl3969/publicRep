package config

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
	Expire int // 小时
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "8.134.63.11",
			Port:     "3306",
			User:     "root",
			Password: "11111",
			Name:     "golang_system",
		},
		JWT: JWTConfig{
			Secret: "123456",
			Expire: 24,
		},
		Server: ServerConfig{
			Port: "8080",
		},
	}
}
