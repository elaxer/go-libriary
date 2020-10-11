package config

// Config ...
type Config struct {
	SecretKey string `toml:"secret_key"`
	DB        struct {
		Engine   string `toml:"engine"`
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
		SSLMode  string `toml:"ssl_mode"`
	} `toml:"db"`
	Server struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"server"`
}
