package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var cats = []categories{
	{ID: 1, Name: "Teori Dasar", Description: "Kumpulan artikel seputar teori dasar elektronika."},
	{ID: 2, Name: "Teknik Digital", Description: "Kumpulan artikel seputar teknik digital."},
	{ID: 3, Name: "Microcontroller", Description: "Kumpulan artikel seputar pemrograman microcontroller."},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
	}

	for _, cat := range cats {
		if cat.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cat)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func updateCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var CatToUpdate categories
	err = json.NewDecoder(r.Body).Decode(&CatToUpdate)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	for i := range cats {
		if cats[i].ID == id {
			CatToUpdate.ID = id
			cats[i] = CatToUpdate

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cats[i])
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i, p := range cats {
		if p.ID == id {
			cats = append(cats[:i], cats[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Category deleted",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

func main() {
	// localhost:2323/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Elektrodukasi V1.0.0",
		})
	})

	// Get all Categories dan POST new kategori
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cats)
		} else if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			var category categories
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			}

			category.ID = len(cats) + 1
			cats = append(cats, category)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(category)
		} else {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	})

	// Get 1 Kategori, update, atau delete
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategoryByID(w, r)
		} else if r.Method == "DELETE" {
			deleteCategoryByID(w, r)
		}
	})

	fmt.Println("Mendengarkan di Localhost pada port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Gagl menjalankan api server")
	}
}
