package datasource

import (
	"gonews/pkg/storage"
)

// Config - данные для обработчика.
type Config struct {
	Links         []string          `json:"rss"`
	RequestPeriod int               `json:"request_period"`
	PostChan      chan storage.Post `json:"-"`
	ErrorChan     chan error        `json:"-"`
}

// Interface задаёт контракт на работу с источником данных.
type Interface interface {
	Run() // запуск источника данных
}
