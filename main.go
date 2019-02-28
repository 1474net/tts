package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/evalphobia/google-tts-go/googletts"
	"github.com/gorilla/mux"
)

var file string
var token string

func tokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func tts(text string) {
	url, err := googletts.GetTTSURL(text, "ru")

	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	resp, _ := http.Get(url)

	defer resp.Body.Close()

	// Create the file
	token = tokenGenerator()
	out, _ := os.Create("./tmp/" + token + ".wav")

	defer out.Close()

	// Write the body to file
	_, _ = io.Copy(out, resp.Body)

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	tts(text)
	fmt.Fprint(w, "/file/"+token+".wav")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/tts/", productsHandler)
	http.Handle("/", router)

	router.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir("./tmp"))))

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
