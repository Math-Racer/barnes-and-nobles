import { useEffect, useState } from "react";
import axios from "axios";
import "./App.css";

function App() {
  const [categories, setCategories] = useState([]);

  useEffect(() => {
    axios.get("https://vijayragav.pythonanywhere.com/categories")
      .then(res => setCategories(res.data))
      .catch(err => console.error(err));
  }, []);

  return (
    <div className="App">
      <h1>Barnes and Nobles</h1>
      {categories.map((cat, idx) => (
        <div key={idx}>
          <h2>{cat.category}</h2>
          <ul>
            {cat.books.map((book, i) => (
              <li key={i}>{book.title} - {book.author}</li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}

export default App;
