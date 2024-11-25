package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// User представляет собой структуру пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

// Инициализация базы данных
func initDB() {
	var err error
	connStr := "user=postgres password=123456 host=localhost port=1138 dbname=Lab8 sslmode=disable" // Замените username и mydb на ваши значения
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

// Обработка ошибок
func handleError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
}

// Валидация данных пользователя
func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return fmt.Errorf("имя не может быть пустым")
	}
	if user.Age < 0 {
		return fmt.Errorf("возраст не может быть отрицательным")
	}
	return nil
}

// Получение списка пользователей с пагинацией и фильтрацией
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Получение параметров запроса
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	nameFilter := r.URL.Query().Get("name")
	ageFilterStr := r.URL.Query().Get("age")

	// Установка значений по умолчанию
	page := 1
	limit := 10

	// Преобразование параметров пагинации
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Формирование SQL-запроса с фильтрацией
	query := "SELECT id, name, age FROM users WHERE TRUE"
	var args []interface{}
	argCount := 1

	if nameFilter != "" {
		query += fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d)", argCount)
		args = append(args, "%"+strings.ToLower(nameFilter)+"%")
		argCount++
	}
	if ageFilterStr != "" {
		ageFilter, err := strconv.Atoi(ageFilterStr)
		if err == nil {
			query += fmt.Sprintf(" AND age = $%d", argCount)
			args = append(args, ageFilter)
			argCount++
		}
	}

	// Добавление пагинации
	query += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			handleError(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

// Получение информации о конкретном пользователе
func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handleError(w, fmt.Errorf("некорректный ID"), http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		handleError(w, fmt.Errorf("пользователь не найден"), http.StatusNotFound)
		return
	} else if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Добавление нового пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validateUser(user); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO users(name, age) VALUES($1, $2) RETURNING id", user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Обновление информации о пользователе
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handleError(w, fmt.Errorf("некорректный ID"), http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	if err := validateUser(updatedUser); err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", updatedUser.Name, updatedUser.Age, id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удаление пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		handleError(w, fmt.Errorf("некорректный ID"), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", r)
}
func TestGetUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users?page=1&limit=2", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUsers)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var users []User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Greater(t, len(users), 0)
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/users/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var user User
	err = json.Unmarshal(rr.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID) // Предполагается, что пользователь с ID 1 существует
}

func TestCreateUser(t *testing.T) {
	newUser := User{Name: "Test User", Age: 30}
	userJSON, _ := json.Marshal(newUser)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdUser User
	err = json.Unmarshal(rr.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Name, createdUser.Name)
	assert.Equal(t, newUser.Age, createdUser.Age)
}

func TestUpdateUser(t *testing.T) {
	updatedUser := User{Name: "Updated User", Age: 35}
	userJSON, _ := json.Marshal(updatedUser)

	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
