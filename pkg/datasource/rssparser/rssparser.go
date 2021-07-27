package rssparser

import (
	"encoding/xml"
	"fmt"
	"gonews/pkg/storage"
	"io"
)

type Parser struct {
	postChan  chan storage.Post
	errorChan chan error
}

//создает объект парсера RSS с заданными параметрами
func New(postChan chan storage.Post, errorChan chan error) *Parser {
	return &Parser{
		postChan:  postChan,
		errorChan: errorChan,
	}
}

//читает RSS
func (s *Parser) Parse(body io.Reader) ([]storage.Post, error) {

	decoder := xml.NewDecoder(body)

	// Чтение item по частям
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.errorChan <- fmt.Errorf("rssparser.Parse_decoder.Token error: %s", err)
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
					s.errorChan <- fmt.Errorf("rssparser.Parse_decoder.DecodeElement error: %s", err)
					return
				}
				s.postChan <- post
			}
		}
	}
}
