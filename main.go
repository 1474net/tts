package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/evalphobia/google-tts-go/googletts"
	"github.com/gorilla/mux"
)

func init() {
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		os.Mkdir("./tmp", os.ModePerm)
	}
}

var token string
var folder string

func tokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func mergefile() string {
	// "-q", "-o 1.mp3", "D:/go/src/tts/tmp/2ce1af05.wav", "D:/go/src/tts/tmp/3aa32875.wav"
	token = tokenGenerator()
	cmd := exec.Command("./mp3cat", "-d", "./tmp/"+folder, "-o", "./tmp/"+folder+"/"+token+".wav")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "1"
	}
	fmt.Printf("The date is %s\n", out)

	return "str"

}
func textleng(text string) {
	var j int
	for len(text) > 200 {

		i := 200
		for text[i] != ' ' {

			i--
		}
		texttmp := text[0:i]

		text = text[i:]
		tts(texttmp, j)
		j++
		fmt.Println(i)
	}
	tts(text, j+1)
}
func tts(text string, i int) {
	url, err := googletts.GetTTSURL(text, "ru")

	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	resp, _ := http.Get(url)

	defer resp.Body.Close()

	// Create the file

	out, _ := os.Create("./tmp/" + folder + "/" + strconv.Itoa(i) + "_" + tokenGenerator() + ".mp3")

	defer out.Close()

	// Write the body to file
	_, _ = io.Copy(out, resp.Body)

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	folder = tokenGenerator()
	if _, err := os.Stat("./tmp/" + folder); os.IsNotExist(err) {
		os.Mkdir("./tmp/"+folder, os.ModePerm)
	}

	text := r.FormValue("text")
	if len(text) > 200 {
		textleng(text)
	} else {
		tts(text, 0)
	}
	mergefile()
	fmt.Fprint(w, "/file/"+folder+"/"+token+".wav")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/tts/", productsHandler)
	http.Handle("/", router)

	router.PathPrefix("/file/").Handler(http.StripPrefix("/file/", http.FileServer(http.Dir("./tmp"))))

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
