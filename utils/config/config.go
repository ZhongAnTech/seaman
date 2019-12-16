package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Sync struct {
		Second time.Duration `yaml:"second"`
	} `yaml:"sync"`
	Git struct {
		Dir    string `yaml:"dir"`
		URL    string `yaml:"url"`
		Branch string `yaml:"branch"`
		Token  string `yaml:"token"`
	} `yaml:"git"`
	Kubecloud struct {
		URL       string `yaml:"url"`
		Token     string `yaml:"token"`
		CompanyId int64  `yaml:"companyId"`
		Cluster   string `yaml:"cluster"`
	} `yaml:"kubecloud"`
}

var (
	config *Config
)

func GetConfig() *Config {
	if config == nil {
		config = &Config{}
		LoadConfig()
	}
	return config
}

var getSearchPath func() []string

func init() {
	getSearchPath = func() []string {
		return []string{
			"configs/config.local.yaml",
			"configs/config.yaml",
		}
	}
}

func LoadConfig() {
	var configPath string
	searchPaths := getSearchPath()
	for _, searchPath := range searchPaths {
		if _, err := os.Stat(searchPath); err != nil {
			continue
		}
		bytes, err := ioutil.ReadFile(searchPath)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(bytes, &config); err != nil {
			panic(err)
		}
		configPath = searchPath
		break
	}
	if configPath == "" {
		panic(fmt.Errorf("config not found: %v", searchPaths))
	}
}
