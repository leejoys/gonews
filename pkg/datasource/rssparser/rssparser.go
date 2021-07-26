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
	links         []string
	requestPeriod int
	postChan      chan storage.Post
	errorChan     chan error
}

//создает объект парсера RSS с заданными параметрами
func New(c datasource.Config) *Source {
	return &Source{
		links:         c.Links,
		requestPeriod: c.RequestPeriod,
		postChan:      c.PostChan,
		errorChan:     c.ErrorChan,
	}
}

//запускает опрос заданных RSS с заданным периодом
func (s *Source) Run() {
	for {
		for _, link := range s.links {
			go s.parse(link)
		}
		time.Sleep(time.Minute * time.Duration(s.requestPeriod))
	}
}

//читает RSS
func (s *Source) parse(link string) {
	resp, err := http.Get(link)
	if err != nil {
		s.errorChan <- fmt.Errorf("rssParse http.Get error: %s", err)
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
			s.errorChan <- fmt.Errorf("rssParse decoder.Token error: %s", err)
			return
		}
		//выбор токена по типу
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "item" {
				// Декодирование элемента в структуру
				var post storage.Post
				err = decoder.DecodeElement(&post, &tp)
				if err != nil {
					s.errorChan <- fmt.Errorf("rssParse decoder.DecodeElement error: %s", err)
					return
				}
				s.postChan <- post
			}
		}
	}
}
