package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Movie struct {
	ID       string    `json:"ID"`
	ISBN     int       `json:"ISBN"`
	Title    string    `json:"Title"`
	Director *Director `json:"Director"`
}

type Director struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updated Movie
	_ = json.NewDecoder(r.Body).Decode(&updated)

	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies[i] = updated
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	var alreadyPresent bool
	_ = json.NewDecoder(r.Body).Decode(&movie)
	for _, v := range movies {
		if v.ID == movie.ID {
	 		alreadyPresent = true
	 		break
		}
	}
	if !alreadyPresent {
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movie)
}

func runHttpServer() {
	r := mux.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	movies = append(movies, Movie{
		ID:       "1",
		ISBN:     123456,
		Title:    "Inception",
		Director: &Director{FirstName: "Chris", LastName: "Nolan"},
	})

	api := r.PathPrefix("/movies").Subrouter()
	api.HandleFunc("", getMovies).Methods("GET")
	api.HandleFunc("", createMovie).Methods("POST")
	api.HandleFunc("/{id}", deleteMovie).Methods("DELETE")
	api.HandleFunc("/{id}", updateMovie).Methods("PUT")

	// Serve static files from ./static
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))

	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
