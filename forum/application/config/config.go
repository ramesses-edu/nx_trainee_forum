package config

import (
	"os"
	"strings"

	"golang.org/x/oauth2"
)

type DBCfg struct {
	UserDB string
	PassDB string
	HostDB string
	PortDB string
	NameDB string
}
type GoogleAuthCfg struct {
	Config *oauth2.Config
}
type FacebookAuthCfg struct {
	Config     *oauth2.Config
	APIVersion string
}
type TwitterAuthCfg struct {
	TwitterAPIKey      string
	TwitterAPISecret   string
	TwitterTokenKey    string
	TwitterTokenSecret string
	RedirectURL        string
	ReqTokenURL        string
	AuthURL            string
	TokenURL           string
}

type Config struct {
	DB       DBCfg
	Google   GoogleAuthCfg
	Facebook FacebookAuthCfg
	Twitter  TwitterAuthCfg
	HostAddr string
}

func New() *Config {
	return &Config{
		DB: DBCfg{
			UserDB: getEnv("USER_DB", ""),
			PassDB: getEnv("PASS_DB", ""),
			HostDB: getEnv("HOST_DB", ""),
			PortDB: getEnv("PORT_DB", ""),
			NameDB: getEnv("NAME_DB", ""),
		},
		Google: GoogleAuthCfg{
			Config: &oauth2.Config{
				ClientID:     getEnv("GA_CLIENT_ID", ""),
				ClientSecret: getEnv("GA_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GA_REDIRECT_URL", ""),
				Scopes:       getEnvAsSlice("GA_SCOPES", []string{}, ","),
				Endpoint: oauth2.Endpoint{
					AuthURL:  getEnv("GA_AUTH_URL", ""),
					TokenURL: getEnv("GA_TOKEN_URL", ""),
				},
			},
		},
		Facebook: FacebookAuthCfg{
			Config: &oauth2.Config{
				ClientID:     getEnv("FBA_CLIENT_ID", ""),
				ClientSecret: getEnv("FBA_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("FBA_REDIRECT_URL", ""),
				Scopes:       getEnvAsSlice("FBA_SCOPES", []string{}, ","),
				Endpoint: oauth2.Endpoint{
					AuthURL:  getEnv("FBA_AUTH_URL", ""),
					TokenURL: getEnv("FBA_TOKEN_URL", ""),
				},
			},
			APIVersion: getEnv("FBA_API_VERSION", "v10.0"),
		},
		Twitter: TwitterAuthCfg{
			TwitterAPIKey:      getEnv("TA_TWITTER_API_KEY", ""),
			TwitterAPISecret:   getEnv("TA_TWITTER_API_SECRET", ""),
			TwitterTokenKey:    getEnv("TA_TWITTER_TOKEN_KEY", ""),
			TwitterTokenSecret: getEnv("TA_TWITTER_TOKEN_SECRET", ""),
			RedirectURL:        getEnv("TA_REDIRECT_URL", ""),
			ReqTokenURL:        getEnv("TA_REQUEST_TOKEN_URL", ""),
			AuthURL:            getEnv("TA_AUTH_URL", ""),
			TokenURL:           getEnv("TA_TOKEN_URL", ""),
		},
		HostAddr: getEnv("HOST_ADDRESS", "localhost:80"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")
	if valStr == "" {
		return defaultVal
	}
	val := strings.Split(valStr, sep)
	return val
}
