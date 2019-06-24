package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Article struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}
type Usuario struct {
	Nombre string `json:"nombre"`
	Edad   string `json:"edad"`
}
type Numeros struct {
	Numero1 int64 `json:"numero1"`
	Numero2 int64 `json:"numero2"`
}
type Respons struct {
	Succes bool   `json:"succes"`
	Code   int64  `json:"code"`
	Result string `json:"result"`
}
type Articles []Article
type Usuarios []Usuario

func allArticles(w http.ResponseWriter, r *http.Request) {

	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Description", Content: "Hello World"},
	}
	fmt.Print("Endpoint Hit: All Articles Endpoint")
	json.NewEncoder(w).Encode(articles)
}

func usuarioInfo(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var usuario Usuario
	err := decoder.Decode(&usuario)
	if err != nil {
		panic(err)
	}
	fmt.Println(usuario.Nombre)
	fmt.Print("Endpoint Hit: Usuario")
	json.NewEncoder(w).Encode(usuario)
}

func proceso(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var numero Numeros
	err := decoder.Decode(&numero)
	if err != nil {
		panic(err)
	}
	suma := numero.Numero1 + numero.Numero2
	resultado := "el resultado es: " + strconv.FormatInt(suma, 10)
	respons := Respons{Succes: true, Code: 200, Result: resultado}
	fmt.Print("Endpoint Hit: Usuario")
	json.NewEncoder(w).Encode(respons)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Dani")
}
func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func handleRequests() {

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", allArticles)
	r.HandleFunc("/recuperarDatos", usuarioInfo).Methods("POST")
	r.HandleFunc("/suma", proceso).Methods("POST")
	r.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	r.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	r.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func baseDatosMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	fmt.Println(client)
}

func main() {
	baseDatosMongoDB()
	handleRequests()
}
