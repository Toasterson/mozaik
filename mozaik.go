package main

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"io"
	"os"
	"log"
	"github.com/toasterson/mozaik/auth"
	"github.com/gorilla/schema"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"fmt"
	"html/template"
)

var schemaDec = schema.NewDecoder()
var pages = map[string]Page{
	"testing": {"Testing", []byte("Testing")},
}

type Pages map[string]Page

type Page struct {
	Title string
	Body []byte
}

func loadPage(title string) (Page, error) {
	return pages[title], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "templates/view", &p)
}


func main() {
	// Initialize Sessions and Cookies
	// Typically gorilla securecookie and sessions packages require
	// highly random secret keys that are not divulged to the public.
	//
	// In this example we use keys generated one time (if these keys ever become
	// compromised the gorilla libraries allow for key rotation, see gorilla docs)
	// The keys are 64-bytes as recommended for HMAC keys as per the gorilla docs.
	//
	// These values MUST be changed for any new project as these keys are already "compromised"
	// as they're in the public domain, if you do not change these your application will have a fairly
	// wide-opened security hole. You can generate your own with the code below, or using whatever method
	// you prefer:
	//
	//    func main() {
	//        fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)))
	//    }
	//
	// We store them in base64 in the example to make it easy if we wanted to move them later to
	// a configuration environment var or file.
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	auth.CookieStore = securecookie.New(cookieStoreKey, nil)
	auth.SetSessionStore(sessions.NewCookieStore(sessionStoreKey))

	// Initialize ab.
	auth.SetupAuthboss()

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)
	router := mux.NewRouter()

	// Routes
	gets := router.Methods("GET").Subrouter()

	router.PathPrefix("/auth").Handler(auth.Ab.NewRouter())

	gets.HandleFunc("/view/{title}", viewHandler)


	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not found")
	})

	// Set up our middleware chain
	stack := alice.New(auth.Logger, auth.Nosurfing, auth.Ab.ExpireMiddleware).Then(router)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}

/*
func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, r, "index", data)
}
*/

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}