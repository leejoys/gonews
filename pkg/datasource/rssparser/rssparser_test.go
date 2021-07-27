package rssparser

import (
	"fmt"
	"gonews/pkg/storage"
	"strings"
	"testing"
	"time"
)

//юнит-тест для Parser_Parse
func TestParser_Parse(t *testing.T) {
	body := strings.NewReader("<rss><channel><title>GOSAMPLES - Learn Golang programming by example</title><link>https://gosamples.dev/</link><description>Learn Golang programming by example. </description><generator>Hugo -- gohugo.io</generator><language>en-us</language><image><url>https://gosamples.dev/apple-touch-icon.png</url><title>GOSAMPLES - Learn Golang programming by example</title><link>https://gosamples.dev/</link></image><lastBuildDate>Mon, 25 Jan 2021 00:00:00 +0000</lastBuildDate><item><title>title1</title><link>link1</link><pubDate>Fri, 23 Jul 2021 00:00:00 +0000</pubDate><guid>https://gosamples.dev/convert-int-to-string/</guid><description>description1</description></item><item><title>title2</title><link>link2</link><pubDate>Tue, 01 Jun 2021 00:00:00 +0000</pubDate><guid>https://gosamples.dev/write-csv/</guid><description>description2</description></item></channel></rss>")
	pc := make(chan storage.Post)
	ec := make(chan error)
	p := &Parser{
		postChan:  pc,
		errorChan: ec,
	}
	go p.Parse(body)
	time.Sleep(time.Second / 4)
	select {
	case err := <-p.errorChan:
		t.Fatalf("err=%s, want nil", err)
	default:
		posts := []storage.Post{}
		for p := range p.postChan {
			posts = append(posts, p)
		}
		fmt.Printf("posts: %v\n", posts)
	}
}
