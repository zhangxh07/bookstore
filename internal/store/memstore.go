package store

import (
	mystore "bookstore/store"
	"bookstore/store/conf"
	factory "bookstore/store/factory"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
)

func init() {
	factory.Register("mem", &MemStore{
		books: make(map[string]*mystore.Book),
	})
}

var (
	ctx = context.Background()
)

type MemStore struct {
	sync.RWMutex
	books map[string]*mystore.Book
	col   *mongo.Collection
}

func (ms *MemStore) Init() error {
	// 补充Mongo依赖
	db, err := conf.C().MongoDB.GetDB()
	if err != nil {
		panic(err)
	}
	ms.col = db.Collection("bookstore")
	return nil
}

// Create creates a new Book in the store.
func (ms *MemStore) Create(book *mystore.Book) error {
	//ms.Lock()
	//defer ms.Unlock()
	//
	//if _, ok := ms.books[book.Id]; ok {
	//	return mystore.ErrExist
	//}
	//
	////nBook := *book
	////ms.books[book.Id] = &nBook
	//ms.books[book.Id] = book

	_, err := ms.col.InsertOne(ctx, book)
	if err != nil {
		return err
	}
	return nil
}

// Update updates the existed Book in the store.
func (ms *MemStore) Update(book *mystore.Book) error {
	_, err := ms.Get(book.Id)
	if err != nil {
		return mystore.ErrNotFound
	}
	res, err := ms.col.UpdateOne(ctx, bson.D{{"id", book.Id}}, bson.D{{"$set", bson.M{
		"name":    book.Name,
		"authors": book.Authors,
		"press":   book.Press,
	}}})

	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

// Get retrieves a book from the store, by id. If no such id exists. an
// error is returned.
func (ms *MemStore) Get(id string) (*mystore.Book, error) {
	//ms.RLock()
	//defer ms.RUnlock()
	//
	//t, ok := ms.books[id]
	//if ok {
	//	return *t, nil
	//}
	//return mystore.Book{}, mystore.ErrNotFound

	ins := &mystore.Book{}
	err := ms.col.FindOne(ctx, bson.M{"id": id}).Decode(ins)
	fmt.Printf("ins: %v\n", ins)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, mystore.ErrNotFound
	}
	return ins, nil
}

// Delete deletes the book with the given id. If no such id exist. an error
// is returned.
func (ms *MemStore) Delete(id string) error {
	//ms.Lock()
	//defer ms.Unlock()
	//
	//if _, ok := ms.books[id]; !ok {
	//	return mystore.ErrNotFound
	//}
	//
	//delete(ms.books, id)
	//return nil
	ins, err := ms.Get(id)
	if err != nil {
		return mystore.ErrNotFound
	}
	res, err := ms.col.DeleteOne(ctx, bson.M{"id": ins.Id})
	if err != nil {
		return err
	}
	log.Printf("delete result: %v\n", res)

	return nil
}

// GetAll returns all the books in the store, in arbitrary order.
func (ms *MemStore) GetAll() ([]mystore.Book, error) {
	//ms.RLock()
	//defer ms.RUnlock()
	//
	//allBooks := make([]mystore.Book, 0, len(ms.books))
	//for _, book := range ms.books {
	//	allBooks = append(allBooks, *book)
	//}
	//return allBooks, nil
	filter := bson.M{}
	res, err := ms.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	set := []mystore.Book{}
	for res.Next(ctx) {
		ins := &mystore.Book{}
		err = res.Decode(ins)
		if err != nil {
			return nil, err
		}
		set = append(set, *ins)
	}
	return set, err
}
