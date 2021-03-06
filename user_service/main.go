package main

import (
	"consul_fabio_demo/utility"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

const (
	fabiopath = "http://172.31.27.241:9999"
)

func main() {
	log.Println("Running User Service...")
	c := utility.GetConsulClient()
	tags := []string{"urlprefix-/service"}
	(*c).Register("UserService", 8890, &tags)

	router := mux.NewRouter()
	router.HandleFunc("/health", health).Methods(http.MethodGet)
	router.HandleFunc("/service/user/{email}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/service/user", createUser).Methods(http.MethodPost)
	router.HandleFunc("/service/user/{email}", deleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/service/user/org/{org}", getOrgUsers).Methods(http.MethodGet)
	router.HandleFunc("/service/users", getUsers).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8890", router))
}

func health(w http.ResponseWriter, r *http.Request) {
	return
}

func getUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Requeset for user info")
	params := mux.Vars(r)
	email := params["email"]

	res, code, err := utility.SendRequest(fabiopath+"/storage/user/"+email, http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostForm.Get("email")
	fname := r.PostForm.Get("fname")
	lname := r.PostForm.Get("lname")
	org := r.PostForm.Get("organization")
	title := r.PostForm.Get("title")
	address := r.PostForm.Get("address")
	log.Printf("Request to create user %s %s %s %s %s %s\n", email, fname, lname, org, title, address)

	data := url.Values{}
	data.Set("email", email)
	data.Set("fname", fname)
	data.Set("lname", lname)
	data.Set("organization", org)
	data.Set("title", title)
	data.Set("address", address)

	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, code, err := utility.SendRequest(fabiopath+"/storage/user", http.MethodPost, strings.NewReader(data.Encode()), &header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to delete user")
	params := mux.Vars(r)
	email := params["email"]

	res, code, err := utility.SendRequest(fabiopath+"/storage/user/"+email, http.MethodDelete, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func getOrgUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get organization users")
	params := mux.Vars(r)
	org := params["org"]

	res, code, err := utility.SendRequest(fabiopath+"/storage/user/org/"+org, http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get all users")
	res, code, err := utility.SendRequest(fabiopath+"/storage/users", http.MethodGet, nil, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(code)
	fmt.Fprint(w, string(*res))
}
