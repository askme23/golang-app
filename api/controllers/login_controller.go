package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
  "log"
	"os"
	"net/http"
  _ "time"

	_ "github.com/askme23/golang-app/api/auth"
	_ "github.com/askme23/golang-app/api/models"
	_ "github.com/askme23/golang-app/api/responses"
	_ "github.com/askme23/golang-app/api/utils/formaterror"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	Email    string `json:"email""`
	Password string `json:"password"`
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Just skip and are waiting request with another method
	if r.Method == "OPTIONS" {
		return
	}  

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		dbname   = os.Getenv("DB_NAME")
		password = os.Getenv("DB_PASSWORD")
	)
	dbConnect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)

	db, err := sql.Open("postgres", dbConnect)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You are Successfully connected!")

  // Checking the existence of the user in the person table.
  var storePassword string
  row := db.QueryRow("select password from person where email = $1", person.Email)
  err = row.Scan(&storePassword)
  if err == sql.ErrNoRows {
    http.Error(w, "Incorrect email or password.", 404)
    return
    // hash, _ := Hash(person.Password)
    // result, err := db.Exec("INSERT INTO person(email, password, when_created, when_updated) VALUES($1, $2, $3, $4)", person.Email, hash, time.Now(), time.Now())

    // if err != nil {
    //     fmt.Println(err)
    //     http.Error(w, http.StatusText(500), 500)
    //     return
    // }

    // if _, err := result.RowsAffected(); err != nil {
    //     http.Error(w, http.StatusText(500), 500)
    //     return
    // }

    // fmt.Println("User is added")
  } else if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if len(storePassword) > 0 {
    _, err := server.SignIn(person.Password, storePassword)

    if err != nil {
      http.Error(w, http.StatusText(500), 500)
      return
    }
  }






	fmt.Fprintf(w, "it's login request")

	// user.Prepare()
	// err = user.Validate("login")
	// if err != nil {
	//  responses.ERROR(w, http.StatusUnprocessableEntity, err)
	//  return
	// }
	// token, err := server.SignIn(user.Email, user.Password)
	// if err != nil {
	//  formattedError := formaterror.FormatError(err.Error())
	//  responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
	//  return
	// }
	// responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(password string, storePassword string) (string, error) {

	var err error

	err = VerifyPassword(password, storePassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
  fmt.Println("Passwords are equal")
	return "", nil//auth.CreateToken(user.ID)
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
