package configs

type Configs struct {
	AppParams      AppParams
	PostgresParams PostgresParams
	AuthParams     AuthParams
	RedisParams    RedisParams
}

type AppParams struct {
	ServerURL  string
	ServerName string
	PortRun    string
	GinMode    string
}

type PostgresParams struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type AuthParams struct {
	AccessTokenTtlMinutes int
	RefreshTokenTtlDays   int
	JwtSecret             string
}

type RedisParams struct {
	Host     string
	Port     string
	Password string
	DB       int
}
