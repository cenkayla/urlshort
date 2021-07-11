package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cenkayla/shorturl/models"
	"github.com/gorilla/mux"
)

type Service struct {
	DB     *models.DB
	router *mux.Router
}

type SaveURLRequestBody struct {
	LongURL   string `json:"long_url"`
	CustomURL string `json:"custom_url"`
}

func InitService(db *models.DB) *Service {
	s := &Service{
		DB:     db,
		router: mux.NewRouter(),
	}
	s.configureRouter()

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) configureRouter() {
	s.router.HandleFunc("/create", s.saveURL).Methods("POST")
	s.router.PathPrefix("/").HandlerFunc(s.redirectByShortURL).Methods("GET")
}

func (s *Service) saveURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := SaveURLRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURL, err := s.DB.SaveURL(data.LongURL, data.CustomURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) redirectByShortURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, err := s.DB.GetLongURL(shortURL)
	if errors.Is(err, models.ErrShortUrlNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
}
