package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Error string
}

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}
func (s *APIServer) Listen() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleAccount))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	log.Printf("JSON server listening on %s", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, router))

}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	case "PATCH":
		return s.handleUpdateAccount(w, r)
	}
	return fmt.Errorf("%s request method not allowed on /account", r.Method)
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)

	if params["id"] != "" {
		account, err := s.store.GetAccountById(params["id"])
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	} else {
		accounts, err := s.store.GetAllAccounts()
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, accounts)
	}

}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountDto := CreateAccountDto{}
	if err := json.NewDecoder(r.Body).Decode(&createAccountDto); err != nil {
		return err
	}
	account, err := s.store.CreateAccount(NewAccount(createAccountDto.FirstName, createAccountDto.LastName))
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
