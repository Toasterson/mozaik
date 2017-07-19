package controller

import (
	"net/http"
	"github.com/toasterson/mozaik/router"
	"gopkg.in/authboss.v1"
)

type ControllerInterface interface {
	Init()			//method = Initializes the Controller
	Prepare(w http.ResponseWriter, r *http.Request) //method = SetUp All Local Variables needed for the Functions
	Get() error     //method = GET processing
	Post() error    //method = POST processing
	Delete() error  //method = DELETE processing
	Put() error     //method = PUT handling
	Head() error    //method = HEAD processing
	Patch() error   //method = PATCH treatment
	Options() error //method = OPTIONS processing
	Finish()		//method = Used to clear Temporary Variables and Cleanup
	GetRoutes()	[]router.Route //Get the Routes of the Controller
	GetAuth() *authboss.Authboss //Get the Authentication Provider
	IsAuthProtected() bool //Return true if Authentication needs to be chacked for this Controller
	AddData(dat map[interface{}]interface{}) //Add Data to the Map available to the templates.
}