package main

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"os"
	"log"
	"github.com/toasterson/mozaik/auth"
	"github.com/gorilla/schema"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/toasterson/mozaik/controller"
	"github.com/justinas/alice"
	"github.com/toasterson/mozaik/xrsf"
)

var schemaDec = schema.NewDecoder()


var (
	mainCont       = MainController{}
	newTileCont    = NewTileController{}
	tileDetailCont = TileDetailController{}
	listTileCont   = TileListController{}
)

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
	//auth.SetupAuthboss()

	//Load Templates to cache
	controller.Init("./templates/*.html", true, "", nil)

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)

	//Routes
	mux_router := mux.NewRouter()

	controller.SetUpRouting(mux_router, []controller.ControllerInterface{
		&mainCont,
		&newTileCont,
		&tileDetailCont,
		&listTileCont,
	})
	mux_router.NotFoundHandler = http.HandlerFunc(controller.NotFound)

	// Set up our middleware chain
	stack := alice.New(auth.Logger, xrsf.Nosurfing).Then(mux_router)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}