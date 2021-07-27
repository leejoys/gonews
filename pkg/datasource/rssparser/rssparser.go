package rssparser

import (
	"encoding/xml"
	"fmt"
	"gonews/pkg/storage"
	"io"
)

type Parser struct{}

//создает объект парсера RSS с заданными параметрами
func New() *Parser {
	return &Parser{}
}

//читает RSS
func (s *Parser) Parse(body io.Reader) ([]storage.Post, error) {

	decoder := xml.NewDecoder(body)
	posts := []storage.Post{}
	// Чтение item по частям
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("rssparser.Parse_decoder.Token error: %s", err)
		}
		//выбор токена по типу
		switch tp := tok.(type) {
		case xml.StartElement:
			if tp.Name.Local == "item" {
				// Декодирование элемента в структуру
				var post storage.Post
				err = decoder.DecodeElement(&post, &tp)
				if err != nil {
					return nil, fmt.Errorf("rssparser.Parse_decoder.DecodeElement error: %s", err)
				}
				posts = append(posts, post)
			}
		}
	}
	return posts, nil
}
