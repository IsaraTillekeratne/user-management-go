package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/user-management/internal/handler"
	"example.com/user-management/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbConn *sql.DB
var userStore *store.UserStore

func TestMain(m *testing.M) {
	ctx := context.Background()

	// temporary postgres container for tests
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "useradmin",
			"POSTGRES_PASSWORD": "userpassword",
			"POSTGRES_DB":       "userdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer postgres.Terminate(ctx)

	host, err := postgres.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=useradmin password=userpassword dbname=userdb sslmode=disable", host, port.Port())
	dbConn, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	schema := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	phone TEXT NOT NULL,
	age INT,
	status TEXT NOT NULL DEFAULT 'Active' CHECK (status IN ('Active','Inactive')),
	created_at TIMESTAMP NOT NULL DEFAULT now()
);
`
	_, err = dbConn.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	userStore = store.NewUserStore(dbConn)

	code := m.Run()
	os.Exit(code)

}

func setupRouter() http.Handler {
	r := chi.NewRouter()
	userHandler := handler.NewUserHandler(userStore)

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetAllUsers)
	r.Get("/users/{id}", userHandler.GetUserById)
	r.Patch("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	return r
}

func TestUserEndpoints(t *testing.T) {
	router := setupRouter()

	// 1. Create User
	userReq := map[string]interface{}{
		"firstName": "Alice",
		"lastName":  "Smith",
		"email":     "alice@example.com",
		"phone":     "+94771234567",
		"age":       28,
		"status":    "Active",
	}

	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", rr.Code, rr.Body.String())
	}

	var createResp map[string]any
	err := json.Unmarshal(rr.Body.Bytes(), &createResp)
	if err != nil {
		t.Fatal(err)
	}

	userData := createResp["user"].(map[string]any)
	userID := userData["UserId"].(string)

	// 2. Get User by ID
	req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	// 3. Update User
	updateReq := map[string]interface{}{
		"firstName": "AliceUpdated",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPatch, "/users/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 on update, got %d", rr.Code)
	}

	// 4. Get All Users
	req = httptest.NewRequest(http.MethodGet, "/users", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 on get all, got %d", rr.Code)
	}

	// 5. Delete User
	req = httptest.NewRequest(http.MethodDelete, "/users/"+userID, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 on delete, got %d", rr.Code)
	}

	// 6. Get Deleted User
	req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleted user, got %d", rr.Code)
	}
}
