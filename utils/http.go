// http
package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	u "net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetHttpClient(userProxy *ProxyS) *http.Client {
	var httpTransport *http.Transport
	if userProxy == nil {

	} else if userProxy.Kind == "sock5" {
		dialer, err := proxy.SOCKS5("tcp", userProxy.Host,
			&proxy.Auth{
				User:     userProxy.User,
				Password: userProxy.Password,
			},
			&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			},
		)
		if err != nil {
			log.Fatalln("get dialer error", dialer)
		}
		httpTransport = &http.Transport{Dial: dialer.Dial}
	} else if userProxy.Kind == "http" {
		proxyUrl, _ := u.Parse(userProxy.Url)
		httpTransport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

	}
	var httpClient *http.Client
	if userProxy == nil {
		httpClient = &http.Client{}
	} else {
		httpClient = &http.Client{Transport: httpTransport}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar = jar
	return httpClient

}

func Get(httpClient *http.Client, url string) ([]byte, error) {
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	reqest.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	reqest.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqest.Header.Set("Cache-Control", "no-cache")

	response, err := httpClient.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return []byte{1}, err
	}

	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return []byte{1}, err
		}
	} else {
		reader = response.Body
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{1}, err
	}

	return body, nil
}

func Post(httpClient *http.Client, url string, v u.Values) ([]byte, error) {
	reqest, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	if err != nil {
		log.Panic(err)
	}
	reqest.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	reqest.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqest.Header.Set("Cache-Control", "no-cache")

	response, err := httpClient.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return []byte{1}, err
	}

	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return []byte{1}, err
		}
	} else {
		reader = response.Body
	}
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{1}, err
	}

	return body, nil
}

func DownloadImage(httpClient *http.Client, url string, fileUrl string) error {
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	reqest.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	reqest.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqest.Header.Set("Cache-Control", "no-cache")

	response, err := httpClient.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		log.Printf("图片下载失败，状态代码: %d ，图片地址: %s", response.StatusCode, url)
		return nil
	}
	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
	} else {
		reader = response.Body
	}
	//open a file for writing
	file, err := os.Create(fileUrl)
	if err != nil {
		return err
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

//jandown.com
func GetJandownTorrent(httpClient *http.Client, url string, fileUrl string) error {
	indexPage, err := Get(httpClient, url)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(indexPage)))
	if err != nil {
		return err
	}
	dom.Find("form").Each(func(i int, s *goquery.Selection) {
		token, _ := s.Find("input[type=text]").Attr("value")
		log.Printf("jandown 获取token: %s", token)
		v := u.Values{}
		v.Add("code", token)
		log.Print(v.Encode())
		reqest, err := http.NewRequest("POST", "http://www.jandown.com/fetch.php", strings.NewReader(v.Encode()))
		if err != nil {
			log.Print(err)
		}
		reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
		reqest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
		reqest.Header.Set("Connection", "keep-alive")
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		reqest.Header.Set("Cache-Control", "no-cache")
		response, err := httpClient.Do(reqest)
		defer response.Body.Close()
		if err != nil {
			log.Print(err.Error())
		}
		if response.StatusCode != 200 {
			log.Printf("种子下载失败，状态代码: %d ，种子地址: %s", response.StatusCode, url)
			log.Print(err.Error())
		}

		var reader io.ReadCloser
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(response.Body)
			if err != nil {
				log.Print(err.Error())
			}
		} else {
			reader = response.Body
		}
		//open a file for writing
		file, err := os.Create(fileUrl)
		if err != nil {
			log.Print(err.Error())
		}
		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, reader)
		if err != nil {
			log.Print(err.Error())
		}
		file.Close()

	})
	return nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
