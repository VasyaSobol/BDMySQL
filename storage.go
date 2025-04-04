package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        string
	Firstname string `validate:"required"`
	Lastname  string `validate:"required"`
	Email     string `validate:"required,email"`
	Age       int    `validate:"required"`
	Created   time.Time
}

type UserInput struct {
	Firstname string `validate:"required" json:"firstname"`
	Lastname  string `validate:"required" json:"lastname"`
	Email     string `validate:"required,email" json:"email"`
	Age       int    `validate:"required" json:"age"`
}

type Server struct {
	Database *sql.DB
}

var validate *validator.Validate

// EditPage godoc
//
//	@Summary      Show a user
//	@Description  get string by ID
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param 	id path string true "ID"
//	@Success      200  {array}   User
//	@Router       /user/{id}	 [get]
func (s *Server) EditPage(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // добавлено для json

	vars := mux.Vars(r)
	ID := vars["id"]

	row := s.Database.QueryRow("select * from userdb.users where ID = ?", ID)

	p := User{}
	var Created string
	err := row.Scan(&p.ID, &p.Firstname, &p.Lastname, &p.Email, &p.Age, &Created)

	if err != nil {
		log.Println(err)
		//	http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		p.Created, err = time.Parse("2006-01-02 15:04:05", Created)
		if err != nil {
			log.Println(err)
			panic(err)
		}
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Println(err)
		panic(err)
	}

}

// EditHandler godoc
//
//	@Summary      updated user ID
//	@Description  updated user ID
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param		  id 	path 	string true "id"
//	@Param        user  body    UserInput  true  "User input"
//	@Router       /edit/{id} [patch]
func (s *Server) EditHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // добавлено для json

	vars := mux.Vars(r)
	ID := vars["id"]

	row := s.Database.QueryRow("select * from userdb.users where ID = ?", ID)
	if row.Err() != nil {
		log.Println("User not found")
		return
	}

	var input UserInput

	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &input)

	// Валидация введенных данных перед обновлением пользователя
	validate = validator.New()
	err := validate.Struct(input)

	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	result, err := s.Database.Exec("update userdb.users set Firstname=?, Lastname=?, Email=?, Age = ? where ID = ?",
		input.Firstname, input.Lastname, input.Email, input.Age, ID)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())

	if err != nil {
		log.Println(err)
	}
}

// CreateHandler godoc
//
//	@Summary      new user
//	@Description  post user
//	@Tags         users
//	@Accept       json
//	@Produce      json
//	@Param        user  body      UserInput  true  "User data"
//	@Success      200  {array}   User
//	@Router       /create [post]
func (s *Server) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var user UserInput

	json.NewDecoder(r.Body).Decode(&user)

	id := uuid.New().String()

	// Валидация введенных данных перед созданием пользователя
	validate = validator.New()
	err := validate.Struct(user)

	if err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusBadRequest, "Validation Error")
		return
	}

	w.Header().Set("Content-Type", "application/json") // добавлено для json
	_, err = s.Database.Exec("insert into userdb.users (ID, Firstname, Lastname, Email, Age, Created) values (?, ?, ?, ?,?,?)",
		id, user.Firstname, user.Lastname, user.Email, user.Age, time.Now())

	if err != nil {
		log.Println(err)
	}
}

// IndexHandler godoc
//
//	@Summary      List users
//	@Description  get users
//	@Tags         users
//	@Accept       json
//	@Produce      json
//
// @Router       / [get]
func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // добавлено для json

	rows, err := s.Database.Query("select * from userdb.users")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {

		p := User{}
		var Created string
		err := rows.Scan(&p.ID, &p.Firstname, &p.Lastname, &p.Email, &p.Age, &Created)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			continue
		}

		p.Created, err = time.Parse("2006-01-02 15:04:05", Created)
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			continue
		}

		users = append(users, p)
	}

	//w.Header().Set("Content-Type", "application/json") // добавлено для json
	json.NewEncoder(w).Encode(users)

	// tmpl, _ := template.ParseFiles("templates/index.html")
	// tmpl.Execute(w, users)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
