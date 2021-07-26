package datasource

import (
	"gonews/pkg/datasource/rssparser"
	"gonews/pkg/storage"
	"time"
)

// Config - данные для обработчика.
type Source struct {
	Links         []string          `json:"rss"`
	RequestPeriod int               `json:"request_period"`
	PostChan      chan storage.Post `json:"-"`
	ErrorChan     chan error        `json:"-"`
	parser        Interface         `json:"-"`
}

// Interface задаёт контракт на работу с источником данных.
type Interface interface {
	Parse(string) // запуск источника данных
}

//запускает опрос заданных адресов с заданным периодом
func (s *Source) Run() {
	p := rssparser.New(s.PostChan, s.ErrorChan)
	s.parser = p
	for {
		for _, link := range s.Links {
			go s.parser.Parse(link)
		}
		time.Sleep(time.Minute * time.Duration(s.RequestPeriod))
	}
}
