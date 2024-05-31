package main

import (
	"log"

	"gorm.io/gorm"
)

type Book struct {
	// ID          int64
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)

	if result.Error != nil {
		// log.Fatalf("Error creating book: %v", result.Error)
		return result.Error
	}

	// fmt.Println("Create Book Successful")
	return nil
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id) // (return value to &book, find record from id)

	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	return &book
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)

	if result.Error != nil {
		log.Fatalf("Error creating book: %v", result.Error)
	}

	return books
}

func updateBook(db *gorm.DB, book *Book) error {
	// result := db.Save(&book) // (return value to &book, find record from id)

	result := db.Model(&book).Updates(book)

	if result.Error != nil {
		// log.Fatalf("Update Book failed: %v", result.Error)
		return result.Error
	}

	// fmt.Println("Update Book Successful")
	return nil
}

func deleteBook(db *gorm.DB, id uint) error {
	var book Book
	result := db.Delete(&book, id) //? Soft Delete
	// result := db.Unscoped().Delete(&book, id) //! Delete permanently

	if result.Error != nil {
		// log.Fatalf("Delete Book failed: %v", result.Error)
		return result.Error
	}

	// fmt.Println("Delete Book Successful")
	return nil
}

func searchBook(db *gorm.DB, bookName string) *Book {
	var book Book

	result := db.Where("name = ?", bookName).First(&book)

	if result.Error != nil {
		log.Fatalf("Search Book failed: %v", result.Error)
	}

	return &book
}

func searchBooks(db *gorm.DB, bookName string) []Book {
	var books []Book

	result := db.Where("name = ?", bookName).Order("price").Find(&books)

	if result.Error != nil {
		log.Fatalf("Search Book failed: %v", result.Error)
	}

	return books
}
