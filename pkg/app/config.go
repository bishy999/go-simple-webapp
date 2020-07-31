package app

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//Conf contains yaml values
type Conf struct {
	SQLConf `yaml:"mysql-config"`
	JWTConf `yaml:"jwt-config"`
}

// SQLConf contains config for initializing the mysql db
type SQLConf struct {
	Name     string `yaml:"name"`
	Usename  string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
}

// JWTConf contains config for jwt
type JWTConf struct {
	Secret string `yaml:"secret"`
}

// AppKey secret key value for JWT
var AppKey string

//GetConf setup database for use
func GetConf() *Conf {

	c := &Conf{}

	yamlFile, err := ioutil.ReadFile("configs/app.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	AppKey = c.Secret
	return c
}
