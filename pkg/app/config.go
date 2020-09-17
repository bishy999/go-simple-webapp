package app

import (
	"io/ioutil"
	"log"
	"os"

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

	// Override with env if set otherwise use yaml
	c.Host = getEnv("MYSQL_HOST", c.Host)
	c.Name = getEnv("MYSQL_DB", c.Name)
	c.Usename = getEnv("MYSQL_USERNAME", c.Usename)
	c.Password = getEnv("MYSQL_PASSWORD", c.Password)
	c.Port = getEnv("MYSQL_PORT", c.Port)

	return c
}


//getEnv find specified env varible
// return value of env variable or if not found return default
func getEnv(key string, defaultVal string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
