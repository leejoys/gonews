package rss

import (
	"encoding/xml"
	"fmt"
	"gonews/pkg/storage"
	"io"
	"log"
	"net/http"
	"time"

	"gonews/pkg/storage/mongodb"
)

//обрабатывает пакет RSS с заданным временем, хранит конфигурацию
type RSSParser struct {
	RSS            []string
	Request_period int
	postChan       chan storage.Post
	errorChan      chan error
	db             *mongodb.Store
}

//запускает опрос заданных RSS
func (p *RSSParser) Run(db *mongodb.Store) {
	p.postChan = make(chan storage.Post)
	p.errorChan = make(chan error)
	p.db = db
	go p.poster()
	go p.logger()
	for {
		for _, link := range p.RSS {
			go p.rssParse(link)
		}
		time.Sleep(time.Minute * time.Duration(p.Request_period))
	}
}

//читает RSS
func (p *RSSParser) rssParse(link string) {
	resp, err := http.Get(link)
	if err != nil {
		p.errorChan <- fmt.Errorf("rssParse http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	// Чтение item по частям
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.errorChan <- fmt.Errorf("rssParse decoder.Token error: %s", err)
			return
		}
		//выбор токена по типу
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "item" {
				// Декодирование элемента в структуру
				var post storage.Post
				decoder.DecodeElement(&post, &tp)
				p.postChan <- post
			}
		}
	}
}

//обрабатывает ответы из каналов с постами
func (p *RSSParser) poster() {
	for post := range p.postChan {

		t, err := time.Parse(time.RFC1123, post.PubDate)
		if err != nil {
			p.errorChan <- fmt.Errorf("poster time.Parse error: %s", err)
		}
		post.PubTime = t.Unix()
		err = p.db.AddPost(post)
		if err != nil {
			p.errorChan <- fmt.Errorf("poster storage.AddPost error: %s", err)
		}
	}
}

//обрабатывает ответы из каналов с ошибками
func (p *RSSParser) logger() {
	for err := range p.errorChan {
		log.Println(err)
	}
}
