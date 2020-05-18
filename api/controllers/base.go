package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/askme23/golang-app/api/models"
	"github.com/askme23/golang-app/api/middlewares"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

var store *sessions.CookieStore
func (server *Server) Initialize(/*Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string*/) {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 60,
		HttpOnly: true,
	}
	// var err error

	// if Dbdriver == "postgres" {
	// 	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	// 	server.DB, err = gorm.Open(Dbdriver, DBURL)
	// 	if err != nil {
	// 		fmt.Printf("Cannot connect to %s database", Dbdriver)
	// 		log.Fatal("This is the error:", err)
	// 	} else {
	// 		fmt.Printf("We are connected to the %s database", Dbdriver)
	// 	}
	// }

	// server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/login", s.Login).Methods("POST", "OPTIONS")
	// s.Router.HandleFunc("/logout", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	// s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	// s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	// s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	// //Posts routes
	// s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatePost)).Methods("POST")
	// s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.GetPosts)).Methods("GET")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetPost)).Methods("GET")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	// s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8081")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}