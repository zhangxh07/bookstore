package store

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrExist    = errors.New("exist")
)

type Book struct {
	Id      string   `json:"id" bson:"id"`           // 图书ISBN ID
	Name    string   `json:"name" bson:"name"`       // 图书名称
	Authors []string `json:"authors" bson:"authors"` // 图书作者
	Press   string   `json:"press" bson:"press"`     // 出版社
}

type Store interface {
	Create(*Book) error
	Update(*Book) error
	Get(string) (*Book, error)
	GetAll() ([]Book, error)
	Delete(string) error
	Init() error
}
