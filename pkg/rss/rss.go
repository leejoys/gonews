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
func (e *RSSParser) Run(db *mongodb.Store) {
	e.postChan = make(chan storage.Post)
	e.errorChan = make(chan error)
	e.db = db
	go e.poster()
	go e.logger()
	for {
		for _, link := range e.RSS {
			go e.rssParse(link)
		}
		time.Sleep(time.Second * time.Duration(e.Request_period))
	}
}

//читает RSS
func (e *RSSParser) rssParse(link string) {
	resp, err := http.Get(link)
	if err != nil {
		e.errorChan <- fmt.Errorf("rssParse http.Get error: %s", err)
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
			e.errorChan <- fmt.Errorf("rssParse decoder.Token error: %s", err)
			return
		}
		//выбор токена по типу
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "item" {
				// Декодирование элемента в структуру
				var p storage.Post
				decoder.DecodeElement(&p, &tp)
				e.postChan <- p
			}
		}
	}
}

//обрабатывает ответы из каналов с постами
func (e *RSSParser) poster() {
	for p := range e.postChan {

		t, err := time.Parse(time.RFC1123, p.PubDate)
		if err != nil {
			e.errorChan <- fmt.Errorf("poster time.Parse error: %s", err)
		}
		p.PubTime = t.Unix()
		err = e.db.AddPost(p)
		if err != nil {
			e.errorChan <- fmt.Errorf("poster storage.AddPost error: %s", err)
		}
	}
}

//обрбатывает ответы из каналов с ошибками
func (e *RSSParser) logger() {
	for err := range e.errorChan {
		log.Println(err)
	}
}
