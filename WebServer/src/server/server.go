/*
TODO: Add a login method that creates and returns a jason web token after successfull login. The key will
need to be generated separately from the jwt-go package.
*/
package main

import (
	//"golang.org/x/crypto/bcrypt"
	"fmt"
	"log"
	"net/http"
	"tools"
)

type App struct {
	UserHandler  *tools.UserHandler
	LoginHandler *tools.LoginHandler
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/users":
		h.UserHandler.ServeHTTP(w, r)
	case "/login":
		h.LoginHandler.ServeHTTP(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func main() {
	/*
		conn := tools.NewConnection()
		err := conn.Connect("mongodb://127.0.0.1:27017")
		if err != nil {
			log.Print(err)
		}
		chore := tools.Chore{
			Name:    "Test Chore",
			Time:    20,
			UserId:  "5df6ade4cc5d561823d8d85e",
			GroupId: "5df6b051cc5d561823d8d860",
		}
		id, err := conn.AddChore(&chore)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("got id %v\n", id)
	*/

	var app App
	tools.Users = map[string]string{
		"Shannon": "password1",
		"Bee":     "password2",
	}
	var userHandler tools.UserHandler
	loginHandler := tools.NewLogin(tools.Users)
	app.UserHandler = &userHandler
	app.LoginHandler = loginHandler
	log.Fatal(http.ListenAndServe(":8080", app))
}
