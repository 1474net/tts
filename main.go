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
	"text/template"

	"github.com/evalphobia/google-tts-go/googletts"
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
	token = tokenGenerator()
	cmd := exec.Command("./mp3/mp3cat", "-d", "./tmp/"+folder, "-o", "./tmp/"+folder+"/"+token+".wav")
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
	if i == -1 {
		token = tokenGenerator()
		out, _ := os.Create("./tmp/" + folder + "/" + token + ".wav")
		defer out.Close()
		_, _ = io.Copy(out, resp.Body)

	} else {
		out, _ := os.Create("./tmp/" + folder + "/" + strconv.Itoa(i) + "_" + tokenGenerator() + ".mp3")
		defer out.Close()
		_, _ = io.Copy(out, resp.Body)

	}

}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	if text == "" {
		fmt.Fprint(w, "Error null text")
		return
	}
	folder = tokenGenerator()
	if _, err := os.Stat("./tmp/" + folder); os.IsNotExist(err) {
		os.Mkdir("./tmp/"+folder, os.ModePerm)
	}

	if len(text) > 200 {
		textleng(text)
		mergefile()
	} else {
		tts(text, -1)
	}

	// http.Redirect(w, r, "/file/"+folder+"/"+token+".wav", http.StatusMovedPermanently)

	http.ServeFile(w, r, "./tmp/"+folder+"/"+token+".wav")
	os.RemoveAll("./tmp/" + folder)

}
func indexHandler(w http.ResponseWriter, r *http.Request) {

	test := struct {
		Test string
	}{
		Test: "test",
	}

	tmpl, _ := template.ParseFiles("templates/index.html")

	tmpl.Execute(w, test)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tts/", productsHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fs := http.FileServer(http.Dir("./tmp"))
	http.Handle("/file/", http.StripPrefix("/file/", fs))

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
