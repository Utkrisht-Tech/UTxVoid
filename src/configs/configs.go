package configs

import (
	// inbuilt go-packages
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
	// created-packages
	"xfiles"
)

// Config settings neccessary for server configuration
type Config struct {
	HTTPHostAndPort  string
	HTTPSHostAndPort string
	HTTPSUsage       string
	HTTPUrl          string
	HTTPSUrl         string
	UseLetsEncrypt   bool
}

// NewConfig generates new configuration settings
func NewConfig() *Config {
	var config Config
	err := config.load()
	if err != nil {
		log.Println("Warning: couldn't load " + xfiles.ConfigFile + ", creating new config file.")
		err = config.create()
		if err != nil {
			log.Fatal("Fatal error: Couldn't create configuration.")
			return nil
		}
		err = config.load()
		if err != nil {
			log.Fatal("Fatal error: Couldn't load configuration.")
			return nil
		}
	}
	return &config
}

// GConfig thread safe global config accessible from all packages
var GConfig = NewConfig()

func (cfg *Config) save() error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(xfiles.ConfigFile, data, 0600)
}

func (cfg *Config) load() error {
	cfgChanged := false
	data, err := ioutil.ReadFile(xfiles.ConfigFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	// Validate and Sanitize HTTPUrl
	if strings.HasSuffix(cfg.HTTPUrl, "/") {
		cfg.HTTPUrl = cfg.HTTPUrl[0 : len(cfg.HTTPUrl)-1]
		cfgChanged = true
	}
	if !strings.HasPrefix(cfg.HTTPUrl, "http://") && !strings.HasPrefix(cfg.HTTPUrl, "https://") {
		cfg.HTTPUrl = "http://" + cfg.HTTPUrl
		cfgChanged = true
	}
	// Validate and Sanitize HTTPSUrl
	if strings.HasSuffix(cfg.HTTPSUrl, "/") {
		cfg.HTTPSUrl = cfg.HTTPSUrl[0 : len(cfg.HTTPSUrl)-1]
		cfgChanged = true
	}
	if strings.HasPrefix(cfg.HTTPSUrl, "http://") {
		cfg.HTTPSUrl = strings.Replace(cfg.HTTPSUrl, "http://", "https://", 1)
		cfgChanged = true
	} else if !strings.HasPrefix(cfg.HTTPSUrl, "https://") {
		cfg.HTTPSUrl = "https://" + cfg.HTTPSUrl
		cfgChanged = true
	}
	// Remove any trailing slashes at the end of the url
	if strings.HasSuffix(cfg.HTTPSUrl, "/") {
		cfg.HTTPSUrl = cfg.HTTPSUrl[0 : len(cfg.HTTPSUrl)-1]
		cfgChanged = true
	}
	// Check if all fields are filled
	cfgReflected := reflect.ValueOf(*cfg)
	for i := 0; i < cfgReflected.NumField(); i++ {
		if cfgReflected.Field(i).Interface() == "" {
			log.Println("Error: " + xfiles.ConfigFile + " is corrupted. Please fill out all of the fields!!")
			return errors.New("Error: configuration corrupted")
		}
	}
	// Save the changed config
	if cfgChanged {
		err = cfg.save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *Config) create() error {
	cfg = &Config{HTTPHostAndPort: ":2120", HTTPSHostAndPort: ":21020", HTTPSUsage: "None", HTTPUrl: "127.0.0.1:2120", HTTPSUrl: "127.0.0.1:21020"}
	err := cfg.save()
	if err != nil {
		log.Println("Error: couldn't create " + xfiles.ConfigFile)
		return err
	}

	return nil
}
