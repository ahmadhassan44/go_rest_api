package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	models "github.com/ahmadhassan44/go_rest_api/models"
	storage "github.com/ahmadhassan44/go_rest_api/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIError struct {
	Error string `json:"error"`
}

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var ae *models.AccountError
			if errors.As(err, &ae) {
				WriteJSON(w, ae.StatusCode, APIError{Error: ae.Error()})
				return
			}
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}
func (s *APIServer) Listen() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/transfer", makeHttpHandlerFunc(s.handleTransfer)).Methods("POST")
	router.HandleFunc("/login", makeHttpHandlerFunc(s.handleLogin)).Methods("POST")
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
	createAccountDto := &models.CreateAccountDto{}
	if err := json.NewDecoder(r.Body).Decode(createAccountDto); err != nil {
		return err
	}
	hashedPasword, err := hashPassword(createAccountDto.Password)
	if err != nil {
		return models.NewAccountError("Could not store you data securely", http.StatusInternalServerError)
	}
	account, err := s.store.CreateAccount(models.NewAccount(createAccountDto.FirstName, createAccountDto.LastName, createAccountDto.UserName, hashedPasword))
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, account)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		return models.NewAccountError("No ID specified for account deletion", http.StatusBadRequest)
	}
	err := s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "Account deleted successfully!")
}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		return models.NewAccountError("No ID specified for account updation", http.StatusBadRequest)
	}
	updateAccountDto := models.UpdateAccountDto{}
	json.NewDecoder(r.Body).Decode(&updateAccountDto)
	err := s.store.UpdateAccount(id, &updateAccountDto)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, "Account updated successfully!")
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferMoneyDto := &models.TransferMoneyDto{}
	json.NewDecoder(r.Body).Decode(transferMoneyDto)
	if transferMoneyDto.Amount == 0 || transferMoneyDto.ReceiverId == "" {
		return models.NewAccountError("Please specify receiver and amount to send!", http.StatusBadRequest)
	}
	err := s.store.TransferMoney(transferMoneyDto)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "Amount Transffered Successfully!")
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginDto := &models.LoginDto{}
	json.NewDecoder(r.Body).Decode(loginDto)
	dbPass, err := s.store.GetUserByUserName(loginDto.UserName)
	if err != nil {
		return err
	}
	if err := verifyPassword(*dbPass, loginDto.Password); err != nil {
		return models.NewAccountError("Invalid Creadentials!", http.StatusUnauthorized)
	}
	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	accessToken, err := generateAccessToken(loginDto.UserName, accessTokenExpiry)
	if (err) != nil {
		return err
	}
	refreshTokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshToken, err := generateAccessToken(loginDto.UserName, refreshTokenExpiry)
	if (err) != nil {
		return err
	}
	err = s.store.SaveRefreshToken(loginDto.UserName, refreshToken, refreshTokenExpiry)
	if (err) != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
func verifyPassword(hashedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
}
func generateAccessToken(userName string, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": userName,          // subject (user ID)
		"exp": expiry.Unix(),     // expiration time
		"iat": time.Now().Unix(), // issued at time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
