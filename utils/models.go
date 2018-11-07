// models.go
package utils

type Topic struct {
	Title string
	Url   string
}

type record struct {
	Id           string
	Title        string
	DownloadTime int64
	IsRequest    bool
	RequestTime  int64
}
