package controllers

import (
  "fmt"
  "net/http"
  _ "io/ioutil"
  "encoding/json"
  
  "github.com/askme23/golang-app/api/auth"
  "github.com/askme23/golang-app/api/models"
  _ "github.com/askme23/golang-app/api/responses"
  _ "github.com/askme23/golang-app/api/utils/formaterror"
  "golang.org/x/crypto/bcrypt"
)

type Person struct {
  Email string `json:"email""`
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

  // body, err := ioutil.ReadAll(r.Body)
  // if err != nil {
  //   http.Error(w, "Something gona wrong", http.StatusUnprocessableEntity)
  //   return
  // }
  
  // fmt.Printf("%s\n", body)
  
  // var person Person
  // err = json.Unmarshal(body, &person)
  // if err != nil {
  //   http.Error(w, "Something gona wrong", http.StatusUnprocessableEntity)
  //   return
  // }

  decoder := json.NewDecoder(r.Body)

  var person Person
  err:= decoder.Decode(&person)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  fmt.Printf("%+v", person)


  fmt.Fprintf(w, "it's login request")
  // body, err := ioutil.ReadAll(r.Body)
  // if err != nil {
  //  responses.ERROR(w, http.StatusUnprocessableEntity, err)
  //  return
  // }
  // user := models.User{}
  // err = json.Unmarshal(body, &user)
  // if err != nil {
  //  responses.ERROR(w, http.StatusUnprocessableEntity, err)
  //  return
  // }

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

func (server *Server) SignIn(email, password string) (string, error) {

  var err error

  user := models.User{}

  err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
  if err != nil {
    return "", err
  }
  err = models.VerifyPassword(user.Password, password)
  if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
    return "", err
  }
  return auth.CreateToken(user.ID)
}