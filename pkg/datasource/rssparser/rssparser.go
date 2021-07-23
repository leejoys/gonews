package rssparser

import (
	"encoding/xml"
	"fmt"
	"gonews/pkg/datasource"
	"gonews/pkg/storage"
	"io"
	"net/http"
	"time"
)

type Source struct {
}

//запускает опрос заданных RSS
func New(c datasource.Config) *Source {
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

//запускает опрос заданных RSS
func (s *Source) Run() {
	for {
		for _, link := range p.RSS {
			go p.rssParse(link)
		}
		time.Sleep(time.Minute * time.Duration(p.Request_period))
	}
}

//читает RSS
func (s *Source) Parse(link string) {
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
