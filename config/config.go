package config

type Config struct {
	Plugins          []string
	AppId            uint64
	AccessToken      string
	Admins           []string
	DatabaseUser     string
	DatabasePassword string
	DatabasePort     int
	DatabaseServer   string
	ServerPort       int
	ScrambleServer   string
	RedisServer      string
	RedisPort        int
	RedisPassword    string
	RedisTable       int
	RedisPoolSize    int
}

var Conf *Config

func init() {
	Conf = &Config{}
}
