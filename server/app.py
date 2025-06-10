from flask import Flask, jsonify
from flask_cors import CORS
from pymongo import MongoClient
import os
from dotenv import load_dotenv 

load_dotenv() 

app = Flask(__name__)
CORS(app)

# MongoDB Atlas connection
MONGO_URI = os.getenv("MONGO_URI")  
client = MongoClient(MONGO_URI)
db = client["barnes-and-nobles"]  # Use your database name
books_collection = db["books"]  # Use your collection name

def group_books_by_category(books):
    categories = {}
    for book in books:
        category = book.get("category", "Uncategorized")
        if category not in categories:
            categories[category] = []
        categories[category].append({
            "title": book.get("title"),
            "author": book.get("author")
        })
    return [
        {"category": cat, "books": books}
        for cat, books in categories.items()
    ]

@app.route("/categories")
def get_categories():
    books = list(books_collection.find({}, {"_id": 0}))  # Exclude MongoDB's _id field
    grouped = group_books_by_category(books)
    return jsonify(grouped)

if __name__ == "__main__":
    app.run(debug=True)