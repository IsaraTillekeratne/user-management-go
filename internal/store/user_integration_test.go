package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"example.com/user-management/internal/model"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbConn *sql.DB
var userStore *UserStore

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

	defer func() {
		if err := postgres.Terminate(ctx); err != nil {
			log.Printf("failed to terminate postgres container: %v", err)
		}
	}()

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

	userStore = NewUserStore(dbConn)

	code := m.Run()
	os.Exit(code)

}

func createTestUser(t *testing.T) model.User {
	t.Helper()
	user := model.User{
		FirstName: "Alice",
		LastName:  "Smith",
		Email:     fmt.Sprintf("alice.%s@example.com", uuid.New().String()),
		Phone:     "+12345678901",
		Age:       30,
		Status:    model.StatusActive,
	}

	created, err := userStore.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	return created
}

func TestCreateUser(t *testing.T) {
	user := model.User{
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     fmt.Sprintf("bob.%s@example.com", uuid.New().String()),
		Phone:     "+12345678902",
		Age:       25,
		Status:    model.StatusActive,
	}

	created, err := userStore.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if created.UserId == uuid.Nil {
		t.Errorf("Expected non-nil UserId")
	}
	if created.Status != model.StatusActive {
		t.Errorf("Expected status Active, got %s", created.Status)
	}
}

func TestCreateUser_DefaultStatus(t *testing.T) {
	user := model.User{
		FirstName: "Bob",
		LastName:  "Johnson",
		Email:     fmt.Sprintf("bob.%s@example.com", uuid.New().String()),
		Phone:     "+12345678902",
		Age:       25,
	}

	created, err := userStore.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if created.UserId == uuid.Nil {
		t.Errorf("Expected non-nil UserId")
	}
	if created.Status != model.StatusActive {
		t.Errorf("Expected status Active, got %s", created.Status)
	}
}

func TestGetAllUsers(t *testing.T) {
	createTestUser(t)
	createTestUser(t)

	users, err := userStore.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}

	if len(users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(users))
	}
}

func TestGetUserById(t *testing.T) {
	user := createTestUser(t)

	got, ok, err := userStore.GetUserById(user.UserId)
	if err != nil {
		t.Fatalf("GetUserById failed: %v", err)
	}
	if !ok {
		t.Fatalf("User not found")
	}
	if got.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, got.Email)
	}
}

func TestUpdateUser(t *testing.T) {
	user := createTestUser(t)

	user.FirstName = "UpdatedName"
	updated, ok, err := userStore.UpdateUser(user, user.UserId)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if !ok {
		t.Fatalf("UpdateUser returned not ok")
	}
	if updated.FirstName != "UpdatedName" {
		t.Errorf("Expected FirstName UpdatedName, got %s", updated.FirstName)
	}
}

func TestDeleteUser(t *testing.T) {
	user := createTestUser(t)

	ok, err := userStore.DeleteUser(user.UserId)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if !ok {
		t.Fatalf("DeleteUser returned not ok")
	}

	_, exists, _ := userStore.GetUserById(user.UserId)
	if exists {
		t.Errorf("Expected user to be deleted")
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	_, ok, err := userStore.GetUserById(uuid.New())
	if err != nil {
		t.Fatalf("GetUserById failed: %v", err)
	}
	if ok {
		t.Errorf("Expected user to not exist")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	user := model.User{
		UserId:    uuid.New(),
		FirstName: "NoOne",
		LastName:  "Nobody",
		Email:     "noone@example.com",
		Phone:     "+12345678999",
		Age:       50,
		Status:    model.StatusInactive,
	}
	_, ok, err := userStore.UpdateUser(user, user.UserId)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if ok {
		t.Errorf("Expected update to return ok=false")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	ok, err := userStore.DeleteUser(uuid.New())
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	if ok {
		t.Errorf("Expected delete to return ok=false")
	}
}
