// filter
package utils

import (
	"github.com/importcjj/sensitive"
)

// 白名单
func LikeFilter(topicList []Topic, keywords []string) []Topic {
	if len(keywords) == 0 {
		return topicList
	}
	filter := sensitive.New()
	filter.AddWord(keywords...)
	var filterdList []Topic
	for _, v := range topicList {
		isFind, _ := filter.FindIn(v.Title)
		if isFind {
			filterdList = append(filterdList, v)
		}
	}
	return filterdList
}

// 黑名单
func HateFilter(topicList []Topic, keywords []string) []Topic {
	if len(keywords) == 0 {
		return topicList
	}
	filter := sensitive.New()
	filter.AddWord(keywords...)
	var filterdList []Topic
	for _, v := range topicList {
		isFind, _ := filter.FindIn(v.Title)
		if !isFind {
			filterdList = append(filterdList, v)
		}
	}
	return filterdList

}

func TitleFilter(title string) string {

	filter := sensitive.New()
	filter.AddWord("\\")
	filter.AddWord("/")
	filter.AddWord(":")
	filter.AddWord("*")
	filter.AddWord("?")
	filter.AddWord("\"")
	filter.AddWord("<")
	filter.AddWord(">")
	filter.AddWord("|")
	filter.AddWord(" ")
	filter.AddWord("→")
	ret := filter.Replace(title, 48)
	return ret

}
