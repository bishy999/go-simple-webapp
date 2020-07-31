package app

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type conf struct {
	SQLConf  `yaml:"mysql-config"`
}


// SQLConf contains config for initializing the mysql db
type SQLConf struct {
	Name     string `yaml:"name"`
	Usename  string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("configs/app.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
