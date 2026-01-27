package handler

import (
	"encoding/json"
	"net/http"

	"example.com/user-management/internal/dto"
	"example.com/user-management/internal/mapper"
	"example.com/user-management/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type UserHandler struct {
	store store.UserStoreInterface
}

func NewUserHandler(store store.UserStoreInterface) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid Request Body"
// @Failure 500 {string} string "Failed to Create User"
// @Router /users [post]
func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	// read the request body and decode it as user
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body!", http.StatusBadRequest)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := mapper.CreateUserRequestToModel(req)
	createdUser, err := handler.store.CreateUser(user)

	if err != nil {
		http.Error(w, "Failed to Create User!", http.StatusInternalServerError)
		return
	}

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

// GetAllUsers godoc
// @Summary Retrieve all users
// @Description Get a list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {string} string "Failed to Retrieve Users"
// @Router /users [get]
func (handler *UserHandler) GetAllUsers(w http.ResponseWriter, _ *http.Request) {
	allUsers, err := handler.store.GetAllUsers()

	if err != nil {
		http.Error(w, "Failed to Retrieve Users!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(allUsers)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieve a user by their UUID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Invalid User Id"
// @Failure 404 {string} string "User Not Found"
// @Failure 500 {string} string "Failed to Retrieve User"
// @Router /users/{id} [get]
func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}
	user, ok, err := handler.store.GetUserById(parsedId)

	if err != nil {
		http.Error(w, "Failed to Retrieve User by Id!", http.StatusInternalServerError)
		return
	}

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

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's information by UUID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid Request Body or User Id"
// @Failure 404 {string} string "User Not Found"
// @Failure 500 {string} string "Failed to Update User"
// @Router /users/{id} [patch]
func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}

	var req dto.UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body!", http.StatusBadRequest)
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok, err := handler.store.GetUserById(parsedId)

	if err != nil {
		http.Error(w, "Failed to Retrieve User by Id!", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "User Not Found!", http.StatusNotFound)
		return
	}

	mapper.ApplyUpdateUserRequest(&user, req)
	updatedUser, _, err := handler.store.UpdateUser(user, parsedId)

	if err != nil {
		http.Error(w, "Failed to Update User!", http.StatusInternalServerError)
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by UUID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid User Id"
// @Failure 404 {string} string "User Not Found"
// @Failure 500 {string} string "Failed to Delete User"
// @Router /users/{id} [delete]
func (handler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid User Id!", http.StatusBadRequest)
		return
	}

	ok, err := handler.store.DeleteUser(parsedId)

	if err != nil {
		http.Error(w, "Failed to Delete User!", http.StatusInternalServerError)
		return
	}

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
