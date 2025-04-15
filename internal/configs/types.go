package configs

type (
	Config struct {
		AppName           string         `mapstructure:"APP_NAME"`
		Version           string         `mapstructure:"VERSION"`
		AppEnv            string         `mapstructure:"APP_ENV"`
		Server            ServerConfig   `mapstructure:",squash"`
		Database          DatabaseConfig `mapstructure:",squash"`
		Redis             RedisConfig    `mapstructure:",squash"`
		Security          SecurityConfig `mapstructure:",squash"`
		Logger            LoggerConfig   `mapstructure:",squash"`
		ModulePermissions []string
	}
)

type (
	// ServerConfig menyimpan konfigurasi server
	ServerConfig struct {
		ServerHost     string   `mapstructure:"APP_HOST"`
		ServerPort     string   `mapstructure:"APP_PORT"`
		ServerEnv      string   `mapstructure:"APP_ENV"`
		AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
	}

	// DatabaseConfig menyimpan konfigurasi database
	DatabaseConfig struct {
		URI      string `mapstructure:"MONGO_URI"`
		Database string `mapstructure:"MONGO_DB"`
		Timeout  int    `mapstructure:"MONGO_TIMEOUT"`
	}

	// RedisConfig menyimpan konfigurasi Redis
	RedisConfig struct {
		Host     string `mapstructure:"REDISHOST"`
		Port     int    `mapstructure:"REDISPORT"`
		Password string `mapstructure:"REDISPASSWORD"`
		DB       int    `mapstructure:"REDIS_DB"`
		PoolSize int    `mapstructure:"POOLSIZE"`
		ConnTTL  int    `mapstructure:"CONNTTL"`
	}

	// SecurityConfig menyimpan konfigurasi keamanan aplikasi
	SecurityConfig struct {
		CheckOrigin            bool   `mapstructure:"ACTIVATE_ORIGIN_VALIDATION"`
		RateLimit              int    `mapstructure:"RATE_LIMIT" envDefault:"60"`
		TrustedPlatform        string `mapstructure:"TRUSTED_PLATFORM"`
		ExpectedHost           string `mapstructure:"EXPECTED_HOST"`
		XFrameOptions          string `mapstructure:"X_FRAME_OPTIONS"`
		ContentSecurity        string `mapstructure:"CONTENT_SECURITY_POLICY"`
		XXSSProtection         string `mapstructure:"X_XSS_PROTECTION"`
		StrictTransport        string `mapstructure:"STRICT_TRANSPORT_SECURITY"`
		ReferrerPolicy         string `mapstructure:"REFERRER_POLICY"`
		XContentTypeOpts       string `mapstructure:"X_CONTENT_TYPE_OPTIONS"`
		PermissionsPolicy      string `mapstructure:"PERMISSIONS_POLICY"`
		JWTSecretKey           string `mapstructure:"JWT_SECRET_KEY"`
		JWTExpired             int    `mapstructure:"JWT_EXPIRED" envDefault:"15"`
		JWTRefreshTokenExpired int    `mapstructure:"JWT_REFRESH_TOKEN_EXPIRED" envDefault:"24"`
		// LimiterInstance        *limiter.Limiter
	}

	// LoggerConfig menyimpan konfigurasi logger
	LoggerConfig struct {
		LogLevel string `mapstructure:"LOG_LEVEL"`
	}
)
