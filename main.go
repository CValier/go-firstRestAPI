package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type User struct {
	ID    int    `json:ID`
	Name  string `json:Name`
	Age   int    `json:Age`
	Email string `json:Email`
}
type allUsers []User

var Users = allUsers{
	{
		ID:    1,
		Name:  "Cristopher",
		Age:   21,
		Email: "cvalier@epa.digital",
	},
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Users)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	newUser := User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	newUser.ID = len(Users) + 1
	Users = append(Users, newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)

}

func deleteUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, user := range Users {
		if user.ID == userID {
			Users = append(Users[:i], Users[i+1:]...)
			fmt.Fprintf(w, "The user with ID %v has been remove succesfully", userID)
		}
	}

}

func putUsers(w http.ResponseWriter, r *http.Request) {

	var updatedUser User

	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please enter valid data")
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i := 0; i < len(Users); i++ {
		if Users[i].ID == userID {
			if updatedUser.Email == "" {
				updatedUser.Email = Users[i].Email
			}
			if updatedUser.Name == "" {
				updatedUser.Name = Users[i].Name
			}
			if updatedUser.Age == 0 {
				updatedUser.Age = Users[i].Age
			}
			Users[i].Name = updatedUser.Name
			Users[i].Age = updatedUser.Age
			Users[i].Email = updatedUser.Email
			fmt.Fprintf(w, "The user with ID %v has been updated succesfully", userID)
		}
	}

	w.Header().Set("Content-Type", "application/json")

}

func menu(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is my REST API :)")
}

func main() {

	http.HandleFunc("/", menu)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/createUsers", postUsers)
	http.HandleFunc("/deleteUsers", deleteUsers)
	http.HandleFunc("/updateUsers", putUsers)
	http.ListenAndServe(":4200", nil)

}
