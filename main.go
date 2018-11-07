// main.go
package main

import (
	"io/ioutil"
	"log"

	"github.com/hardseedGO/aisex"
	"github.com/hardseedGO/chaoliu"
	"github.com/hardseedGO/xp"
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
	XPAddr      string   `yaml:"xp_addr,omitempty"`
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

	//init chaoliu struct
	chaoliu.C.Url = c.CaoliuAddr
	chaoliu.C.Timeout = c.Timeout
	chaoliu.C.HateKeyWord = c.HateKeyWord
	chaoliu.C.LikeKeyWord = c.LikeKeyWord
	chaoliu.C.SavePath = c.SavePath
	chaoliu.C.TopicRange = c.TopicRange
	chaoliu.C.Proxy = c.Proxy

	//init xp struct
	xp.C.Url = c.XPAddr
	xp.C.Timeout = c.Timeout
	xp.C.HateKeyWord = c.HateKeyWord
	xp.C.LikeKeyWord = c.LikeKeyWord
	xp.C.SavePath = c.SavePath
	xp.C.TopicRange = c.TopicRange
	xp.C.Proxy = c.Proxy

	for _, AvClass := range c.AVClass {
		switch AvClass {
		case "aicheng_asia_mosaiched":
			aisex.Do("aicheng_asia_mosaiched")
		case "aicheng_asia_non_mosaicked":
			aisex.Do("aicheng_asia_non_mosaicked")
		case "chaoliu_asia_mosaiched":
			chaoliu.Do("chaoliu_asia_mosaiched")
		case "chaoliu_asia_non_mosaiched":
			chaoliu.Do("chaoliu_asia_non_mosaiched")
		case "xp_asia_mosaiched":
			xp.Do("xp_asia_mosaiched")
		case "xp_asia_non_mosaiched":
			xp.Do("xp_asia_non_mosaiched")
		default:
			log.Println("no av class find")
		}
	}

}
