package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text_test_generator/splitter"
	"time"
)

type Texts struct {
	Test     string
	Validate string
}

const separators = "\n .,!?'\"-:;—(){}[]*/\\1234567890"

func getCheck(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("./texts")

	fIndex := rand.Int() % len(files)
	text, _ := ioutil.ReadFile("./texts/" + files[fIndex].Name())

	origText := string(text)
	spl := splitter.Split(origText, separators)
	words := spl.Words()
	wordsCount := len(words) / 5
	proceeded := make(map[int]bool)
	for wordsCount > 0 {
		i := rand.Int() % len(words)
		word := words[i]
		if proceeded[i] || len(*word) < 2 {
			continue
		}
		newWord := []byte(*word)
		for j := len(newWord) / 2; j < len(newWord); j++ {
			newWord[j] = '*'
		}
		*word = string(newWord)
		proceeded[i] = true
		wordsCount--
	}

	testText := spl.Reconfigure()

	tmpl := template.Must(template.ParseFiles("check.html"))
	_ = tmpl.Execute(w, Texts{
		Test:     testText,
		Validate: origText,
	})
}

func getAddText(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("add_text.html"))
	_ = tmpl.Execute(w, nil)
}

func postAddText(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	f, _ := os.Create("./texts/" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt")
	_, _ = f.Write(body)
	defer f.Close()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getCheck).Methods(http.MethodGet)
	r.HandleFunc("/add", getAddText).Methods(http.MethodGet)
	r.HandleFunc("/add", postAddText).Methods(http.MethodPost)

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(httpServer.ListenAndServe())
}
