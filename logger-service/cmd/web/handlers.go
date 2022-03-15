package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"reflect"
	"time"
)

const authServiceURL = "http://authentication-service/authenticate"

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	_ = readJSON(w, r, &requestPayload)

	collection := app.Mongo.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      requestPayload.Name,
		Data:      requestPayload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println(err)
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Message = "logged"
	_ = writeJSON(w, http.StatusAccepted, resp)
}

// LoginPage displays the login page
func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	render(w, "login.page.gohtml")
}

// LoginPagePost handles user login. Note that it calls the authentication microservice
func (app *Config) LoginPagePost(w http.ResponseWriter, r *http.Request) {
	// it's always good to regenerate the session on login/logout
	_ = app.Session.RenewToken(r.Context())

	// parse the posted form data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// create a variable we'll post to the auth service as JSON
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	requestPayload.Email = email
	requestPayload.Password = password

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(requestPayload, "", "\t")

	// call the authentication-service
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("wrong status code")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("did not get status accepted")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user struct {
		ID        int       `json:"id"`
		Email     string    `json:"email"`
		FirstName string    `json:"first_name,omitempty"`
		LastName  string    `json:"last_name,omitempty"`
		Password  string    `json:"-"`
		Active    int       `json:"active"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	_ = readJSON(w, r, &user)

	// set up session & log user in
	app.Session.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

type Logs struct {
	ID        string
	Name      string
	Data      string
	CreatedAt time.Time
}

// Dashboard displays the dashboard page
func (app *Config) Dashboard(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
	}
	defer cursor.Close(ctx)

	log.Println("Got cursor")

	for cursor.Next(ctx) {
		var result bson.M
		var item Logs
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("cursor.Next() error:", err)
			app.clientError(w, http.StatusBadRequest)
		} else {
			log.Println("\nresult type:", reflect.TypeOf(result))
			log.Println("result:", item.Name)
		}
	}

	render(w, "dashboard.page.gohtml")
}