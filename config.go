package figo

type AppConfig struct {
	Port    int    `port for http listen`
	Address string `ip address for http listen . unuse`
	LogPath string `log path`

	CacheAddress  string `cache address`
	CachePassword string `cache password`
	CacheDB       string `cache db`

	DbType     string
	DbHost     string
	DbUser     string
	DbPassword string
	DbName     string
	DbPort     int
}

type AppEnv struct {
	IsDebug bool

	Extra map[string]string
}

var FigoEnv = newAppEnv()

func newAppEnv() *AppEnv {
	re := &AppEnv{}
	re.Extra = make(map[string]string, 16)

	return re
}
