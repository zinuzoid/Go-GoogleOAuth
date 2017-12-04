package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates *template.Template

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Print("recover", r)
			return
		}
		log.Print("recover nil")
	}()

	templates = template.Must(template.ParseFiles("view/gplus.html", "view/header.html", "view/footer.html"))

	fmt.Print("hi!")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login/gplus", loginGplusHandler)
	http.HandleFunc("/login/gplus/verify", loginGplusVerifyHandler)
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func loginGplusHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "gplus.html", struct {
		Title  string
		Header string
	}{"Leet!", "Hello Google OAuth 2.0"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print(r.RequestURI)
}

func loginGplusVerifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	email := r.FormValue("email")

	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var objmap map[string]string
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if objmap["email"] == email {
		fmt.Fprint(w, "true")
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "false")
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404", http.StatusNotFound)
}
