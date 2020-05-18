package controllers

import (
	"net/http"
	"fmt"
	_ "github.com/askme23/golang-app/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	// session, err := store.Get(r, "session-name")
 //  if err != nil {
 //      http.Error(w, err.Error(), http.StatusInternalServerError)
 //      return
 //  }

 	fmt.Println("its home page")
	fmt.Fprintf(w, "Welcome To This Awesome API")
	// responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}