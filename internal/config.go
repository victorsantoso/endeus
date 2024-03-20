package internal

// Config.
type Config struct {
	Database    Database
	Application Application
}

// Application configuration.
type Application struct {
	Port  string
	Debug bool
}

// Database configuration.
type Database struct {
	ConnectionPool   ConnectionPool
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
	DatabaseHost     string
	DatabasePort     string
	DatabaseSslMode  string
}

// Database connection pool.
type ConnectionPool struct {
	MaxOpenConns   int
	MaxLifetime    int
	MaxIdleConns   int
	MaxIdleTime    int
	DefaultTimeout int
}

type JWT struct {
	Secret            string
	Issuer            string
	Audience          string
	Subject           string
	JwtExpirationTime int
}

// Configure Application configuration with spf13/viper
func ConfigureApplication() *Application {
	return &Application{
		Port:  ViperReader.GetString("application.port"),
		Debug: ViperReader.GetBool("application.debug"),
	}
}

// Configure Dsn from config using spf13/viper
func ConfigureDatabase() *Database {
	return &Database{
		ConnectionPool: ConnectionPool{
			MaxOpenConns: ViperReader.GetInt("database.connection_pool.max_open"),
			MaxLifetime:  ViperReader.GetInt("database.connection_pool.max_lifetime"),
			MaxIdleConns: ViperReader.GetInt("database.connection_pool.max_idle"),
			MaxIdleTime:  ViperReader.GetInt("database.connection_pool.max_idle_time"),
		},
		DatabaseUsername: ViperReader.GetString("database.username"),
		DatabasePassword: ViperReader.GetString("database.password"),
		DatabaseName:     ViperReader.GetString("database.name"),
		DatabaseHost:     ViperReader.GetString("database.host"),
		DatabasePort:     ViperReader.GetString("database.port"),
		DatabaseSslMode:  ViperReader.GetString("database.ssl_mode"),
	}
}

// Configure JWT configuration with spf13/viper
func ConfigureJWT() *JWT {
	return &JWT{
		Issuer:            ViperReader.GetString("jwt.issuer"),
		Audience:          ViperReader.GetString("jwt.audience"),
		Secret:            ViperReader.GetString("jwt.secret"),
		Subject:           ViperReader.GetString("jwt.subject"),
		JwtExpirationTime: ViperReader.GetInt("jwt.jwt_expiration_time"),
	}
}
