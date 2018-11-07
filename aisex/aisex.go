package aisex

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hardseedGO/utils"
)

var (
	C config
)

type config struct {
	Url         string
	AVClass     string
	Timeout     int
	TopicRange  string
	SavePath    string
	HateKeyWord []string
	LikeKeyWord []string
	Proxy       string
}

func Do(avClass string) {
	proxy := utils.GetProxy(C.Proxy)
	httpclient := utils.GetHttpClient(proxy)

	var topicList []utils.Topic
	// 计算要请求多少页面
	page, _ := strconv.Atoi(C.TopicRange)
	pageCount := page / 28
	if pageCount == 0 {
		pageCount++
	}

	if (page % 28) != 0 {
		pageCount++
	}
	log.Printf("一共需要请求 %d 页", pageCount)
	//get index page

	for i := 0; i < pageCount; i++ {
		url := getTopicsListUrl(avClass) + "&page=" + strconv.Itoa(i+1)
		log.Println(url)
		body, err := utils.Get(httpclient, url)
		if err != nil {
			log.Panic(err)
		}
		utfBody, err := utils.GbkToUtf8(body)
		if err != nil {
			log.Panic(err)
		}
		tmpTopicList := getTopic(utfBody)
		for _, value := range tmpTopicList {
			topicList = append(topicList, value)
		}

		// 随机休息
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(10)
		time.Sleep(time.Duration(x) * time.Second)

	}
	log.Printf("一共得到title : %d", len(topicList))

	//过滤喜欢
	topicList = utils.LikeFilter(topicList, C.LikeKeyWord)
	log.Printf("过滤喜好主题后主题数 : %d", len(topicList))
	for _, value := range topicList {
		log.Printf("topic: %s - %s\n", value.Title, C.Url+value.Url)
	}

	//过滤不喜欢
	topicList = utils.HateFilter(topicList, C.HateKeyWord)
	log.Printf("过滤不喜好主题后主题数 : %d", len(topicList))
	for _, value := range topicList {
		log.Printf("topic: %s - %s\n", value.Title, C.Url+value.Url)
	}

	//处理过滤后的主题，开始请求图片和种子
	for _, value := range topicList {
		//todo 增加本地状态查询

		//请求详情页面
		topicBody, err := utils.Get(httpclient, C.Url+value.Url)
		if err != nil {
			log.Panic(err)
		}
		utfBody, err := utils.GbkToUtf8(topicBody)
		if err != nil {
			log.Panic(err)
		}

		err = getImageAndTorrent(utfBody, value.Title)
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(time.Duration(1) * time.Second)

	}

}

func getTopicsListUrl(avClass string) string {
	switch avClass {
	case "aicheng_mosaiched":
		return C.Url + "thread.php?fid=4"
	case "aicheng_asia_non_mosaicked":
		return C.Url + "thread.php?fid=16"
	default:
		return C.Url + "thread.php?fid=4"
	}
}

func getTopic(body []byte) []utils.Topic {
	var topicList []utils.Topic
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find(".tr3.t_one").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a[target=_blank]").Text()
		url, _ := s.Find("a[target=_blank]").Attr("href")
		tmpTopic := utils.Topic{
			Title: utils.TitleFilter(title),
			Url:   url,
		}
		topicList = append(topicList, tmpTopic)
	})
	return topicList
}

func getImageAndTorrent(body []byte, title string) error {
	proxy := utils.GetProxy(C.Proxy)
	httpclient := utils.GetHttpClient(proxy)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	dom.Find(".tpc_content").Each(func(i int, s *goquery.Selection) {
		//得到图片
		id, _ := s.Attr("id")
		if id == "read_tpc" {
			s.Find("img").Each(func(i1 int, s1 *goquery.Selection) {
				imageUrl, _ := s1.Attr("src")
				log.Printf("正在下载图片: %s", imageUrl)
				savePath := C.SavePath + "/" + title + "-" + strconv.Itoa(i1) + ".jpg"
				utils.DownloadImage(httpclient, imageUrl, savePath)
			})
			//得到种子地址
			torrent, _ := s.Find("a[target=_blank]").Attr("href")

			log.Printf("获取种子地址： %s", torrent)
			savePath := C.SavePath + "/" + title + ".torrent"

			if strings.Contains(torrent, "rmdown") {
				err = utils.GetRmdownTorrent(httpclient, torrent, savePath)
				if err != nil {
					log.Print(err)
				}
			} else if strings.Contains(torrent, "jandown") {
				err = utils.GetJandownTorrent(httpclient, torrent, savePath)
				if err != nil {
					log.Print(err)
				}
			} else {
				log.Print("不能解析torrent下载方式，请联系开发人员")
			}

		}

	})
	return nil
}
