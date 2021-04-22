package conf

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Google struct {
	Application Application `yaml:"application"`
}

type Application struct {
	Credentials Credentials `yaml:"credentials"`
}

type Credentials struct {
	Storage string `yaml:"storage"`
}

type Conf struct {
	Google Google `yaml:"google"`
}

func ReadConfig(path string) Conf {
	c := Conf{}
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("Cannot find configuration file: " + path)
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Cannot open configuration file: " + path)
	}
	yaml.Unmarshal(buf, &c)
	return c
}
