package rss

import (
	"encoding/xml"
	"fmt"
	"gonews/pkg/storage"
	"io"
	"log"
	"net/http"
	"time"
)

//обрабатывает пакет RSS с заданным временем, хранит конфигурацию
type RSSEater struct {
	RSS            []string
	Request_period int
	postChan       chan storage.Post
	errorChan      chan error
}

//запускает опрос заданных RSS
func (e *RSSEater) Run() {
	e.postChan = make(chan storage.Post)
	e.errorChan = make(chan error)
	go e.poster()
	go e.logger()
	for {
		for _, link := range e.RSS {
			go e.rssEat(link)
		}
		time.Sleep(time.Second * time.Duration(e.Request_period))
	}
}

//читает RSS
func (e *RSSEater) rssEat(link string) {
	resp, err := http.Get(link)
	if err != nil {
		e.errorChan <- fmt.Errorf("rssEat http.Get error: %s", err)
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
			e.errorChan <- fmt.Errorf("rssEat decoder.Token error: %s", err)
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
func (e *RSSEater) poster() {
	for p := range e.postChan {

		t, err := time.Parse(time.RFC1123, p.PubDate)
		if err != nil {
			e.errorChan <- fmt.Errorf("poster time.Parse error: %s", err)
		}

		p.PubTime = t.Unix()
		fmt.Println(p.PubTime)
	}
}

//обрбатывает ответы из каналов с ошибками
func (e *RSSEater) logger() {
	for err := range e.errorChan {
		log.Println(err)
	}
}
