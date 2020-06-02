package controllers

import (
  "net/http"
  "fmt"
  "os"
  _ "github.com/askme23/golang-app/api/responses"
  jwt "github.com/dgrijalva/jwt-go"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
  w.Header().Set("Access-Control-Allow-Credentials", "true")
  w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

  // Just skip and are waiting request with another method
  if r.Method == "OPTIONS" {
    return
  } 

  access, err := r.Cookie("jwt_access")
  fmt.Println(access.Value)
  if err != nil {
  //   if err == http.ErrNoCookie {
  //     refresh, err := r.Cookie("jwt_refresh")
  //     if err != nil && err == http.ErrNoCookie {
  //       w.WriteHeader(http.StatusUnauthorized)
  //       return
  //     }

  //     token, err := VerifyToken(refresh.Value, "refresh")
  //     if err != nil {
  //       return
  //     }
  //     fmt.Println(token)
  //     if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
  //       return
  //     }
		// } 

    // w.WriteHeader(http.StatusBadRequest)
    // return
  }

	// TODO Здесь проверка валидности access токена
  token, err := VerifyToken(access.Value, "access")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
    return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
    w.WriteHeader(http.StatusUnauthorized)
    return
	}

	w.Write([]byte(fmt.Sprintf("%s", "hello Ruslan")))
  
  


  // if err != nil {
  //  if err == http.ErrNoCookie {
  //    w.WriteHeader(http.StatusUnauthorized)
  //    return
  //  }
  //  // For any other type of error, return a bad request status
  //  w.WriteHeader(http.StatusBadRequest)
  //  return
  // }
  // sessionToken := c.Value
  // fmt.Println(sessionToken)

  // We then get the name of the user from our cache, where we set the session token
  // response, err := cache.Do("GET", sessionToken)
  // if err != nil {
  //  // If there is an error fetching from cache, return an internal server error status
  //  w.WriteHeader(http.StatusInternalServerError)
  //  return
  // }
  // if response == nil {
  //  // If the session token is not present in cache, return an unauthorized error
  //  w.WriteHeader(http.StatusUnauthorized)
  //  return
  // }
  // // Finally, return the welcome message to the user
  // w.Write([]byte(fmt.Sprintf("%s", response)))
}

func VerifyToken(tokenString string, tokenType string) (*jwt.Token, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
     }

    if tokenType == "refresh" { 
      return []byte(os.Getenv("REFRESH_SECRET")), nil
    } else {
      return []byte(os.Getenv("ACCESS_SECRET")), nil
    }
  })
  
  if err != nil {
		fmt.Println(err)
		return nil, err
  }

  return token, nil
}