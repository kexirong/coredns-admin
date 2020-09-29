package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const filePath = "config.yaml"

//Config all config in this
type Config struct {
	Host         string `yaml:"host,omitempty"`
	Port         string `yaml:"port,omitempty"`
	UserEtcdPath string `yaml:"user_etcd_path,omitempty"`
	Etcd         struct {
		Endpoint   []string `yaml:"endpoint"`
		PathPrefix string   `yaml:"path_prefix"`
		Username   string   `yaml:"username"`
		Password   string   `yaml:"password"`
		Timeout    int      `yaml:"timeout"`
		TLS        []string `yaml:"tls"`
	}
}

var conf *Config

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func loadConfig(path string) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}

//Get  get a congfig instance
func Get() *Config {
	if conf != nil {
		return conf
	}
	conf = new(Config)

	if pathIsExist(filePath) {
		loadConfig(filePath)
	}

	loadDefaultConfig(conf)

	return conf
}

// func init() {
// 	Get()
// }
