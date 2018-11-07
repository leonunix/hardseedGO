package xp

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
	pageCount := page / 50
	if pageCount == 0 {
		pageCount++
	}

	if (page % 60) != 0 {
		pageCount++
	}
	log.Printf("一共需要请求 %d 页", pageCount)
	//get index page

	for i := 0; i < pageCount; i++ {
		url := getTopicsListUrl(avClass) + "-page-" + strconv.Itoa(i+1) + ".html"
		log.Println(url)
		body, err := utils.Get(httpclient, url)
		if err != nil {
			log.Panic(err)
		}

		tmpTopicList := getTopic(body)
		for _, value := range tmpTopicList {
			topicList = append(topicList, value)
		}

		// 随机休息
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(10)
		time.Sleep(time.Duration(x) * time.Second)

	}
	log.Printf("一共得到title : %d", len(topicList))
	for _, value := range topicList {
		log.Printf("topic: %s - %s\n", value.Title, C.Url+value.Url)
	}
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
	case "xp_asia_mosaiched":
		return C.Url + "thread-htm-fid-22"
	case "xp_asia_non_mosaiched":
		return C.Url + "thread-htm-fid-5"
	default:
		return C.Url + "thread-htm-fid-22"
	}
}

func getTopic(body []byte) []utils.Topic {
	var topicList []utils.Topic
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("h3").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url, _ := s.Find("a").Attr("href")
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
	isFirst := true

	dom.Find("div[class=f14][id=read_tpc]").Each(func(i int, s *goquery.Selection) {

		if isFirst {
			isFirst = false
			s.Find("img").Each(func(i1 int, s1 *goquery.Selection) {
				imageUrl, _ := s1.Attr("src")
				if strings.HasSuffix(imageUrl, ".gif") {
					return
				}
				log.Printf("正在下载图片: %s", imageUrl)
				savePath := C.SavePath + "/" + title + "-" + strconv.Itoa(i1) + ".jpg"
				utils.DownloadImage(httpclient, imageUrl, savePath)
			})
			//得到种子地址
			var torrent string
			s.Find("a[target=_blank]").Each(func(i int, s1 *goquery.Selection) {

				tmpUrl := s1.Text()
				if strings.Contains(tmpUrl, "downsx") {
					torrent = tmpUrl
				}
			})
			log.Printf("获取种子地址： %s", torrent)
			savePath := C.SavePath + "/" + title + ".torrent"
			err = utils.GetDownsxTorrent(httpclient, torrent, savePath)
			if err != nil {
				log.Print(err)
			}

		}

	})
	return nil
}
