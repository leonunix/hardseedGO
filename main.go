// main.go
package main

import (
	"io/ioutil"
	"log"

	"github.com/hardseedGO/aisex"
	"gopkg.in/yaml.v2"
)

const (
	softName = "hardseedGO"
	version  = "0.1"
	author   = "leon.unix"
)

type config struct {
	AisexAddr   string   `yaml:"aisex_addr,omitempty"`
	CaoliuAddr  string   `yaml:"chaoliu_addr,omitempty"`
	AVClass     []string `yaml:"av_class,omitempty"`
	Timeout     int      `yaml:"timeout,omitempty"`
	TopicRange  string   `yaml:"topic_range,omitempty"`
	SavePath    string   `yaml:"save_path,omitempty"`
	HateKeyWord []string `yaml:"hate_keywords,omitempty"`
	LikeKeyWord []string `yaml:"like_keywords,omitempty"`
	Proxy       string   `yaml:"proxy_addr,omitempty"`
}

var (
	c config
)

func (c *config) getConf() *config {
	yamlFile, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c

}
func main() {
	log.Println(softName)
	log.Println(version)
	log.Println(author)
	c.getConf()

	//init aisex struct
	aisex.C.Url = c.AisexAddr
	aisex.C.Timeout = c.Timeout
	aisex.C.HateKeyWord = c.HateKeyWord
	aisex.C.LikeKeyWord = c.LikeKeyWord
	aisex.C.SavePath = c.SavePath
	aisex.C.TopicRange = c.TopicRange
	aisex.C.Proxy = c.Proxy

	for _, AvClass := range c.AVClass {
		switch AvClass {
		case "aicheng_mosaiched":
			aisex.Do("aicheng_mosaiched")
		case "aicheng_asia_non_mosaicked":
			aisex.Do("aicheng_asia_non_mosaicked")
		default:
			log.Println("no av class find")
		}
	}

}
