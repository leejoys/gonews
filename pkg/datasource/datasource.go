package datasource

import (
	"fmt"
	"gonews/pkg/storage"
	"io"
	"net/http"
	"time"
)

// Config - данные для обработчика.
type Source struct {
	Links         []string          `json:"rss"`
	RequestPeriod int               `json:"request_period"`
	PostChan      chan storage.Post `json:"-"`
	ErrorChan     chan error        `json:"-"`
	Parser        Interface         `json:"-"`
}

// Interface задаёт контракт на работу с источником данных.
type Interface interface {
	Parse(io.Reader) // запуск источника данных
}

//запускает опрос заданных адресов с заданным периодом
func (s *Source) Run() {
	for {
		for _, link := range s.Links {
			go func(link string) {
				resp, err := http.Get(link)
				if err != nil {
					s.ErrorChan <- fmt.Errorf("datasource.Run_http.Get error: %s", err)
					return
				}
				defer resp.Body.Close()
				s.Parser.Parse(resp.Body)
			}(link)
		}
		time.Sleep(time.Minute * time.Duration(s.RequestPeriod))
	}
}
