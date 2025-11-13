package config

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	MaxConns int32
	MinConns int32
}

// LoadDBConfigs loads the database configuration
func LoadDBConfigs() *DBConfig {
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "easyqueue"),
		Password: getEnv("DB_PASSWORD", "easyqueue123"),
		Database: getEnv("DB_NAME", "easyqueue"),
		MaxConns: 25,
		MinConns: 5,
	}
}
