package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var filePath = "config.yaml"

//Config all config in this
type Config struct {
	Listen   string `yaml:"listen"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`

	Etcd struct {
		Endpoint   []string `yaml:"endpoint"`
		PathPrefix string   `yaml:"path_prefix"`
		Username   string   `yaml:"username"`
		Password   string   `yaml:"password"`
		Timeout    int      `yaml:"timeout"`
		TLS        []string `yaml:"tls"`
	} `yaml:"etcd"`
	Redis struct {
		Addresses      []string `yaml:"addresses"`
		KeyPrefix      string   `yaml:"key_prefix"`
		Username       string   `yaml:"username"`
		Password       string   `yaml:"password"`
		ConnectTimeout int      `yaml:"connect_timeout"`
		ReadTimeout    int      `yaml:"read_timeout"`
		TLS            []string `yaml:"tls"`
	} `yaml:"redis"`
}

var conf *Config

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func loadConfig(path string) *Config {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	conf = new(Config)
	if err := yaml.Unmarshal([]byte(data), conf); err != nil {
		log.Fatalf("error: %v", err)
	}
	return conf
}

//Get  get a congfig instance
func Get() *Config {
	if conf == nil {
		Set(filePath)
	}
	return conf
}

func Set(filePath string) {

	if !pathIsExist(filePath) {
		log.Fatalln("the configuration file does not exist")
	}

	loadConfig(filePath)
	loadDefaultConfig(conf)

}

func LoadDefaultConfig() {
	conf = new(Config)
	loadDefaultConfig(conf)
}

// func init() {
// 	Get()
// }
