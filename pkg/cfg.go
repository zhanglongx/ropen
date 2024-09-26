package pkg

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	CFG_DEFAULT_PORT = 38080

	CFG_DEFAULT_FILENAME  = APP_NAME + ".yml"
	CFG_DEFAULT_FILENAME2 = APP_NAME + ".yaml"
)

type Config struct {
	Port int

	PreferIPs []string

	Issuer struct {
		CAPath  string
		KeyPath string
	}
}

var Cfg Config = Config{
	Port:      CFG_DEFAULT_PORT,
	PreferIPs: []string{},
}

func LoadCfg(path string) {
	var candidates []string

	if isFileExist(path) {
		candidates = []string{path}
	} else {
		candidates = []string{
			CFG_DEFAULT_FILENAME,
			CFG_DEFAULT_FILENAME2,
		}

		// cfg from HOME directory
		if homeDir, err := os.UserHomeDir(); err == nil {
			candidates = append(candidates,
				filepath.Join(homeDir, "."+CFG_DEFAULT_FILENAME),
				filepath.Join(homeDir, "."+CFG_DEFAULT_FILENAME2),
			)
		}
	}

	for _, c := range candidates {
		if err := loadConfigFromFile(c); err == nil {
			debug("loaded config from %v", c)
			return
		}
	}

	debug("failed to load config from %v, using embedded config",
		candidates)
}

func isFileExist(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func loadConfigFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		return err
	}

	return nil
}
