package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func handleRequests() {

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/articles", allArticles)
	r.HandleFunc("/recuperarDatos", usuarioInfo).Methods("POST")
	r.HandleFunc("/suma", proceso).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func main() {
	handleRequests()
}
