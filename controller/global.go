package controller

import (
	"github.com/dannyvankooten/grender"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/toasterson/mozaik/logger"
	"fmt"
)

var (
	grend *grender.Grender
)

//Initialize Global Server stuff
func Init(tplGlob string, debug bool, partialGlob string, funcs template.FuncMap) {
	grend = grender.New(grender.Options{
		Debug: debug,       // If true, templates will be recompiled before each render call
		TemplatesGlob: tplGlob,  // Glob to your template files
		PartialsGlob: partialGlob,   // Glob to your patials or global templates
		Funcs: funcs,         // Your template FuncMap
		Charset: "UTF-8",   // Charset to use for Content-Type header values
	})
	SetupAuthboss("./templates/auth/*.html", debug, partialGlob, Funcs)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	grend.HTML(w, http.StatusForbidden, "forbidden.html", nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	grend.HTML(w, http.StatusNotFound, "notfound.html", nil)
}

func MakeHandler(controllerInterface ControllerInterface) *Handler{
	return &Handler{controllerInterface}
}

func SetUpRouting(mux_router *mux.Router, controllers []ControllerInterface){
	for _, controller := range controllers {
		logger.Trace(fmt.Sprintf("Setting Up Controller: %s", getName(controller)))
		controller.Init()
		for _, route := range controller.GetRoutes() {
			logger.Trace("Setting up Route:", route.Path, route.Method)
			mux_router.Handle(route.Path, MakeHandler(controller)).Methods(route.Method)
		}
		if controller.IsAuthProtected(){
			logger.Trace(fmt.Sprintf("Enabling Authentication for %s", getName(controller)))
			enableAuth(mux_router)
		}
	}
}