package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/user-management/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MockUserStore struct {
	CreateUserFn  func(model.User) (model.User, error)
	GetAllUsersFn func() ([]model.User, error)
	GetUserByIdFn func(uuid.UUID) (model.User, bool, error)
	UpdateUserFn  func(model.User, uuid.UUID) (model.User, bool, error)
	DeleteUserFn  func(uuid.UUID) (bool, error)
}

func (m *MockUserStore) CreateUser(u model.User) (model.User, error) {
	return m.CreateUserFn(u)
}
func (m *MockUserStore) GetAllUsers() ([]model.User, error) {
	return m.GetAllUsersFn()
}
func (m *MockUserStore) GetUserById(id uuid.UUID) (model.User, bool, error) {
	return m.GetUserByIdFn(id)
}
func (m *MockUserStore) UpdateUser(u model.User, id uuid.UUID) (model.User, bool, error) {
	return m.UpdateUserFn(u, id)
}
func (m *MockUserStore) DeleteUser(id uuid.UUID) (bool, error) {
	return m.DeleteUserFn(id)
}

// unit tests for CreateUser
func TestCreateUser_Success(t *testing.T) {
	mockUserStore := &MockUserStore{
		CreateUserFn: func(user model.User) (model.User, error) {
			user.UserId = uuid.New()
			return user, nil
		},
	}

	userHandler := NewUserHandler(mockUserStore)

	body := `{
		"firstName":"John",
		"lastName":"Doe",
		"email":"john@gmail.com",
		"phone":"+94712345678",
		"age": 27,
		"status": "Active"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	userHandler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	handler := NewUserHandler(&MockUserStore{})

	body := `{
		"firstName":"J",
		"email":"invalid-email",
		"phone":"0712345678"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// unit tests for GetAllUsers
func TestGetAllUsers_Success(t *testing.T) {
	mockUserStore := &MockUserStore{
		GetAllUsersFn: func() ([]model.User, error) {
			return []model.User{
				{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@gmail.com",
					Phone:     "+94712345678",
					Age:       27,
					Status:    "Active",
				},
			}, nil
		},
	}

	userHandler := NewUserHandler(mockUserStore)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	userHandler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// unit tests for GetUserById
func TestGetUserById_Success(t *testing.T) {
	id := uuid.New()

	mockStore := &MockUserStore{
		GetUserByIdFn: func(uid uuid.UUID) (model.User, bool, error) {
			return model.User{
				UserId:    uid,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@gmail.com",
				Phone:     "+94712345678",
				Age:       27,
				Status:    "Active",
			}, true, nil
		},
	}

	handler := NewUserHandler(mockStore)

	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserById)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetUserById_InvalidUUID(t *testing.T) {
	handler := NewUserHandler(&MockUserStore{})

	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserById)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid-id", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	mockStore := &MockUserStore{
		GetUserByIdFn: func(uuid.UUID) (model.User, bool, error) {
			return model.User{}, false, nil
		},
	}

	handler := NewUserHandler(mockStore)
	id := uuid.New()

	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserById)

	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// unit tests for UpdateUser
func TestUpdateUser_Success(t *testing.T) {
	id := uuid.New()

	mockStore := &MockUserStore{
		GetUserByIdFn: func(uuid.UUID) (model.User, bool, error) {
			return model.User{
				UserId:    id,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@gmail.com",
				Phone:     "+94712345678",
				Age:       27,
				Status:    "Active",
			}, true, nil
		},
		UpdateUserFn: func(u model.User, id uuid.UUID) (model.User, bool, error) {
			return u, true, nil
		},
	}

	handler := NewUserHandler(mockStore)

	body := `{"firstName":"Updated"}`
	r := chi.NewRouter()
	r.Patch("/users/{id}", handler.UpdateUser)

	req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateUser_ValidationError(t *testing.T) {
	id := uuid.New()
	handler := NewUserHandler(&MockUserStore{})

	body := `{"email":"bad-email"}`

	req := httptest.NewRequest(http.MethodPatch, "/users/"+id.String(), bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/users/{id}", handler.UpdateUser)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	id := uuid.New()

	mockStore := &MockUserStore{
		GetUserByIdFn: func(uuid.UUID) (model.User, bool, error) {
			return model.User{}, false, nil
		},
	}

	handler := NewUserHandler(mockStore)

	body := `{"firstName":"Jane"}`
	req := httptest.NewRequest("PATCH", "/users/"+id.String(), bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Patch("/users/{id}", handler.UpdateUser)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// unit tests for DeleteUser
func TestDeleteUser_Success(t *testing.T) {
	id := uuid.New()

	mockStore := &MockUserStore{
		DeleteUserFn: func(uuid.UUID) (bool, error) {
			return true, nil
		},
	}

	handler := NewUserHandler(mockStore)

	r := chi.NewRouter()
	r.Delete("/users/{id}", handler.DeleteUser)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteUser_InvalidUUID(t *testing.T) {
	handler := NewUserHandler(&MockUserStore{})

	req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/users/{id}", handler.DeleteUser)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockStore := &MockUserStore{
		DeleteUserFn: func(uuid.UUID) (bool, error) {
			return false, nil
		},
	}

	handler := NewUserHandler(mockStore)
	id := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Delete("/users/{id}", handler.DeleteUser)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
