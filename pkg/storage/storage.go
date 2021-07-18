package storage

// Post - публикация.
type Post struct {
	ID      int    // номер записи
	Title   string `xml:"title" json:"title"`             // заголовок публикации
	Content string `xml:"description" json:"description"` // содержание публикации
	PubDate string `xml:"pubDate" json:"-"`               // время публикации из RSS
	PubTime int64  `xml:"-" json:"pubtime"`               //время публикации для БД и фронта
	Link    string `xml:"link" json:"link"`               // ссылка на источник
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error     // создание новой публикации
	UpdatePost(Post) error  // обновление публикации
	DeletePost(Post) error  // удаление публикации по ID
	Close()                 // освобождение ресурса
}
