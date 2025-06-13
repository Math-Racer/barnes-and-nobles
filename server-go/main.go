package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Title    string `bson:"title" json:"title"`
	Author   string `bson:"author" json:"author"`
	Category string `bson:"category" json:"category"`
}

type CategoryGroup struct {
	Category string `json:"category"`
	Books    []Book `json:"books"`
}

var booksCollection *mongo.Collection

func groupBooksByCategory(books []Book) []CategoryGroup {
	categories := make(map[string][]Book)
	for _, book := range books {
		category := book.Category
		if strings.TrimSpace(category) == "" {
			category = "Uncategorized"
		}
		categories[category] = append(categories[category], Book{
			Title:  book.Title,
			Author: book.Author,
		})
	}
	var grouped []CategoryGroup
	for cat, books := range categories {
		grouped = append(grouped, CategoryGroup{Category: cat, Books: books})
	}
	return grouped
}

func getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var books []Book
	if err := cursor.All(ctx, &books); err != nil {
		http.Error(w, "Error reading books", http.StatusInternalServerError)
		return
	}

	grouped := groupBooksByCategory(books)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grouped)
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("barnes-and-nobles")
	booksCollection = db.Collection("books")

	http.HandleFunc("/categories", getCategoriesHandler)

	log.Println("Server running on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
