package storage

// Post - публикация.
type Post struct {
	ID      int    // номер записи
	Title   string `xml:"title"`       // заголовок публикации
	Content string `xml:"description"` // содержание публикации
	PubDate string `xml:"pubDate"`     // время публикации
	PubTime int64
	Link    string `xml:"link"` // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
	Close()                 // освобождение ресурса
}
