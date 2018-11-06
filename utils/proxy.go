// proxy
package utils

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

type ProxyS struct {
	Kind     string
	Addr     string
	Port     int
	User     string
	Password string
	Host     string
	Url      string
}

func GetProxy(proxy string) *ProxyS {
	if proxy == "" {
		return nil
	}
	u, err := url.Parse(proxy)
	if err != nil {
		log.Fatal(err)
	}
	var userProxy *ProxyS = new(ProxyS)
	userProxy.Url = proxy
	userProxy.Kind = u.Scheme
	userProxy.User = u.User.Username()
	userProxy.Password, _ = u.User.Password()
	h := strings.Split(u.Host, ":")
	userProxy.Addr = h[0]
	port, err := strconv.Atoi(h[1])
	userProxy.Port = port
	userProxy.Host = u.Host
	return userProxy

}
