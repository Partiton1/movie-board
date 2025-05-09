package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	movies = []Movie{} // reset for each test

	r := mux.NewRouter()
	api := r.PathPrefix("/movies").Subrouter()
	api.HandleFunc("", getMovies).Methods("GET")
	api.HandleFunc("", createMovie).Methods("POST")
	api.HandleFunc("/{id}", deleteMovie).Methods("DELETE")
	api.HandleFunc("/{id}", updateMovie).Methods("PUT")
	return r
}

func TestCreateAndGetMovies(t *testing.T) {
	r := setupRouter()

	movie := Movie{
		ID:    "1",
		ISBN:  123456,
		Title: "Test Movie",
		Director: &Director{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	body, _ := json.Marshal(movie)

	req := httptest.NewRequest("POST", "/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Now test GET
	req = httptest.NewRequest("GET", "/movies", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("GET /movies returned status %d", resp.Code)
	}
}

func TestUpdateMovie(t *testing.T) {
	r := setupRouter()

	// Insert a movie first
	movies = append(movies, Movie{
		ID:    "2",
		ISBN:  654321,
		Title: "Old Title",
		Director: &Director{
			FirstName: "Jane",
			LastName:  "Smith",
		},
	})

	updatedMovie := Movie{
		ID:    "2",
		ISBN:  654321,
		Title: "New Title",
		Director: &Director{
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}
	body, _ := json.Marshal(updatedMovie)

	req := httptest.NewRequest("PUT", "/movies/2", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

func TestDeleteMovie(t *testing.T) {
	r := setupRouter()

	// Insert a movie
	movies = append(movies, Movie{
		ID:    "3",
		ISBN:  789012,
		Title: "Delete Me",
		Director: &Director{
			FirstName: "Jim",
			LastName:  "Beam",
		},
	})

	req := httptest.NewRequest("DELETE", "/movies/3", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "3"})
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}
