package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

var (
	Db                   *sql.DB
	serverPort           string
	mySecretStringForJWT string
)

func StartServer(db *sql.DB, servPort, JWTsecret string) error {

	if db != nil {
		Db = db
	} else {
		return errors.New("db is nil! ")
	}

	serverPort = servPort
	mySecretStringForJWT = JWTsecret

	r := chi.NewRouter()
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/deals", func(r chi.Router) {
		r.Post("/create", auth(dealCreate))
		r.Get("/list", auth(dealsList))
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/info", auth(dealInfo))
		})
	})
	fmt.Println("server is listening on ", serverPort)
	err := http.ListenAndServe(":"+serverPort, r)
	log.Fatal("unexpected exit with err: ", err)
	return err
}

func auth(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			resp := Response{
				Error: "No token!",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		var mySigningKey = []byte(mySecretStringForJWT)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing ")
			}
			return mySigningKey, nil
		})

		if err != nil {
			fmt.Println("token expired")
			json.NewEncoder(w).Encode("token expired")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				us, ok := claims["user"]
				if ok {
					user, ok := us.(string)
					if ok {
						r.Header.Set("User", user)
					}
				}
				r.Header.Set("Role", "admin")

				handler.ServeHTTP(w, r)
				return

			}
			// if need more roles
			/*else if claims["role"] == "user" {

				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			}*/
		}
		fmt.Println("not authorized")
		w.WriteHeader(http.StatusUnauthorized)
		resp := Response{
			Error: "Not authorized",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
}
