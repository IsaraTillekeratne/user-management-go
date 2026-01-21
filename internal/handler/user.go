package handler

import (
	"encoding/json"
	"net/http"

	"example.com/user-management/internal/model"
	"example.com/user-management/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(store *store.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	// read the request body and decode it as user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Request Body!", http.StatusBadRequest)
		return
	}

	err = validate.Struct(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser := handler.store.CreateUser(user)

	response := map[string]interface{}{
		"message": "User created successfully!",
		"user":    createdUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers := handler.store.GetAllUsers()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(allUsers)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}
	user, ok := handler.store.GetUserById(parsedId)
	if !ok {
		http.Error(w, "User Not Found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

// UpdateUser TODO: add validator for update request
func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}

	var user model.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Request Body!", http.StatusBadRequest)
		return
	}

	updatedUser, ok := handler.store.UpdateUser(user, parsedId)
	if !ok {
		http.Error(w, "User Not Found!", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "User Updated successfully!",
		"user":    updatedUser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}

	ok := handler.store.DeleteUser(parsedId)
	if !ok {
		http.Error(w, "User Not Found!", http.StatusNotFound)
		return
	}

	response := map[string]string{
		"message": "User Deleted successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
