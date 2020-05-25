package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
  "log"
	"os"
	"net/http"
  "time"

	_ "github.com/askme23/golang-app/api/auth"
	_ "github.com/askme23/golang-app/api/models"
	_ "github.com/askme23/golang-app/api/responses"
	_ "github.com/askme23/golang-app/api/utils/formaterror"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
)

type Person struct {
	Email    string `json:"email""`
	Password string `json:"password"`
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	// Just skip and are waiting request with another method
	if r.Method == "OPTIONS" {
		return
	}  

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var person Person
	var err = json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = godotenv.Load(); err != nil {
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

    // create cookie
    sessionValue := map[string]string {
    	"email": person.Email,
    }

    if hash, err := s.Encode("cookie-login", sessionValue); err == nil {
    	fmt.Println()
    	_, err = cache.Do("SETEX", hash, "120", person.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

      http.SetCookie(w, &http.Cookie{
				Name:    "cookie-login",
				Value:   hash,
				Expires: time.Now().Add(120 * time.Second),
			})
    }

    generateTokens()
  }
}

// В дальнейшем вынести в отдельный файл
func generateTokens() {
	generateAccessToken()
	generateRefreshToken()
}

func generateAccessToken() error {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["name"] = "Gorelov Ruslan"
  claims["admin"] = true
  // claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
  claims["exp"] = time.Now().Add(time.Second * 120).Unix()

  t, err := accessToken.SignedString([]byte(os.Getenv("API_SECRET")))
  if err != nil {
      return err
  }
  fmt.Println(t)
  return nil
}

func generateRefreshToken() error {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["sub"] = 1
  claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

  t, err := refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
  if err != nil {
      return err
  }
  fmt.Println(t)
  return nil
}


func (server *Server) SignIn(password string, storePassword string) (string, error) {
	var err error
	err = VerifyPassword(password, storePassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
  fmt.Println("Passwords are equal")
	return "", nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
