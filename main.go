package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Movie Model
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director Model
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Movies Database
var movies []Movie

/*
* Name: Get all movies
* Method: GET
* Route: /movies
 */
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

/*
* Name: Get Movie By ID
* Method: GET
* Route: /movies/{id}
* Params: ID of the movie
 */
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	movieId := params["id"]

	for _, movie := range movies {
		if movieId == movie.ID {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	http.Error(w, "Movie Not Found", http.StatusNotFound)
}

/*
* API: Create Movie
* Method: POST
* Route: /movies
 */
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	movie.ID = strconv.Itoa(rand.Intn(100000000))

	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movies)
}

/*
* Name: Update Movie By ID
* Method: PUT
* Route: /movies/{id}
* Params: ID of the movie
 */
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	movieId := params["id"]

	movieIndex := -1

	for index, movie := range movies {
		if movieId == movie.ID {
			movieIndex = index
			break
		}
	}

	if movieIndex == -1 {
		http.Error(w, "Movie Not Found", http.StatusNotFound)
		return
	}

	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	movie.ID = movieId
	movies[movieIndex] = movie

	json.NewEncoder(w).Encode(movies)
}

/*
* Name: Delete Movie By ID
* Method: DELETE
* Route: /movies/{id}
* Params: ID of the movie
 */
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	movieId := params["id"]

	movieIndex := -1

	for index, movie := range movies {
		if movieId == movie.ID {
			movieIndex = index
			break
		}
	}

	if movieIndex == -1 {
		http.Error(w, "Movie Not Found", http.StatusNotFound)
		return
	}

	movies = append(movies[:movieIndex], movies[movieIndex+1:]...)
	json.NewEncoder(w).Encode(movies)
}

func main() {
	router := mux.NewRouter()

	movies = append(
		movies,
		Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}},
		Movie{ID: "2", Isbn: "453218", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}
