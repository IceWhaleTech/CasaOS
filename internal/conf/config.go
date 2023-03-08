package conf

type Database struct {
	Type        string `json:"type" env:"DB_TYPE"`
	Host        string `json:"host" env:"DB_HOST"`
	Port        int    `json:"port" env:"DB_PORT"`
	User        string `json:"user" env:"DB_USER"`
	Password    string `json:"password" env:"DB_PASS"`
	Name        string `json:"name" env:"DB_NAME"`
	DBFile      string `json:"db_file" env:"DB_FILE"`
	TablePrefix string `json:"table_prefix" env:"DB_TABLE_PREFIX"`
	SSLMode     string `json:"ssl_mode" env:"DB_SSL_MODE"`
}

type Scheme struct {
	Https    bool   `json:"https" env:"HTTPS"`
	CertFile string `json:"cert_file" env:"CERT_FILE"`
	KeyFile  string `json:"key_file" env:"KEY_FILE"`
}

type LogConfig struct {
	Enable     bool   `json:"enable" env:"LOG_ENABLE"`
	Name       string `json:"name" env:"LOG_NAME"`
	MaxSize    int    `json:"max_size" env:"MAX_SIZE"`
	MaxBackups int    `json:"max_backups" env:"MAX_BACKUPS"`
	MaxAge     int    `json:"max_age" env:"MAX_AGE"`
	Compress   bool   `json:"compress" env:"COMPRESS"`
}

type Config struct {
	Force          bool      `json:"force" env:"FORCE"`
	Address        string    `json:"address" env:"ADDR"`
	Port           int       `json:"port" env:"PORT"`
	SiteURL        string    `json:"site_url" env:"SITE_URL"`
	Cdn            string    `json:"cdn" env:"CDN"`
	JwtSecret      string    `json:"jwt_secret" env:"JWT_SECRET"`
	TokenExpiresIn int       `json:"token_expires_in" env:"TOKEN_EXPIRES_IN"`
	Database       Database  `json:"database"`
	Scheme         Scheme    `json:"scheme"`
	TempDir        string    `json:"temp_dir" env:"TEMP_DIR"`
	BleveDir       string    `json:"bleve_dir" env:"BLEVE_DIR"`
	Log            LogConfig `json:"log"`
}
