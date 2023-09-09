package api

import (
	"encoding/json"
	"go-rest-api/pkg/db/models"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
)

// start api with the pgdb and return a chi router
func StartAPI(pgdb *pg.DB) *chi.Mux {
	//get the router
	r := chi.NewRouter()
	//add middleware
	//in this case we will store our DB to use it later
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))

	r.Route("/cards", func(r chi.Router) {
		r.Post("/", createComment)
		r.Get("/", getComments)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})

	return r
}

type CreateCardRequest struct {
	Card string `json:"card"`
}
type CardResponse struct {
	// Success bool         `json:"success"`
	// Error   string       `json:"error"`
	Card *models.Card `json:"card"`
}

type CardsResponse struct {
	// Success bool           `json:"success"`
	// Error   string         `json:"error"`
	Cards []*models.Card `json:""`
}

func createComment(w http.ResponseWriter, r *http.Request) {
	//get the request body and decode it
	req := &CreateCardRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	//if there's an error with decoding the information
	//send a response with an error
	if err != nil {
		res := &CardResponse{
			// Success: false,
			// Error:   err.Error(),
			Card: nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get the db from context
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	//if we can't get the db let's handle the error
	//and send an adequate response
	if !ok {
		res := &CardResponse{
			// Success: false,
			// Error:   "could not get the DB from context",
			Card: nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if we can get the db then
	card, err := models.CreateCard(pgdb, &models.Card{})
	if err != nil {
		res := &CardResponse{
			// Success: false,
			// Error:   err.Error(),
			Card: nil,
		}
		err = json.NewEncoder(w).Encode(res)
		//if there's an error with encoding handle it
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		//return a bad request and exist the function
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//everything is good
	//let's return a positive response
	res := &CardResponse{
		// Success: true,
		// Error:   "",
		Card: card,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding after creating comment %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getComments(w http.ResponseWriter, r *http.Request) {
	//get db from ctx
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		res := &CardsResponse{
			// Success: false,
			// Error:   "could not get DB from context",
			Cards: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//call models package to access the database and return the comments
	cards, err := models.GetAllCards(pgdb)
	if err != nil {
		res := &CardsResponse{
			// Success: false,
			// Error:   err.Error(),
			Cards: nil,
		}
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Printf("error sending response %v\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//positive response
	res := &CardsResponse{
		// Success: true,
		// Error:   "",
		Cards: cards,
	}
	//encode the positive response to json and send it back
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error encoding comments: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
