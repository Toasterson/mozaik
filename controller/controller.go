package controller

import (
	"html/template"
	"net/http"
	"errors"
	"github.com/gorilla/mux"
	"github.com/toasterson/mozaik/router"
	"mime/multipart"
	"github.com/dannyvankooten/grender"
	"github.com/toasterson/mozaik/logger"
)

var (
	grend *grender.Grender
)

type Controller struct {
	postFormParsed bool
	isInitialized bool
	Routes   []router.Route
	Data     map[interface{}]interface{}
	Name     string
	TplName  string
	TplExt   string
	W        http.ResponseWriter
	R        *http.Request
}

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
}

type Handler struct {
	ControllerInterface
}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.Prepare(w, r)
	defer this.Finish()
	var err error = nil
	switch r.Method {
	case http.MethodGet:
		err = this.Get()
	case http.MethodPost:
		err = this.Post()
	case http.MethodPatch:
		err = this.Patch()
	case http.MethodPut:
		err = this.Put()
	case http.MethodHead:
		err = this.Head()
	case http.MethodDelete:
		err = this.Delete()
	case http.MethodOptions:
		err = this.Options()
	default:
		err = errors.New("No Method")
	}
	//TODO Better Handling Should an error Occur Maybe Custom error Class?
	if err != nil {
		logger.Error(err)
		NotFound(w, r)
	}
}

func MakeHandler(controllerInterface ControllerInterface) *Handler{
	return &Handler{controllerInterface}
}

func SetUpRouting(mux_router *mux.Router, controllers []ControllerInterface){
	for _, controller := range controllers {
		controller.Init()
		for _, route := range controller.GetRoutes() {
			mux_router.Handle(route.Path, MakeHandler(controller)).Methods(route.Method)
		}
	}
}

func (this *Controller) GetRoutes()	[]router.Route{
	return this.Routes
}

func (this *Controller) Prepare(w http.ResponseWriter, r *http.Request){
	this.W = w
	this.R = r
}

func (this *Controller) Finish(){
	this.W = nil
	this.R = nil
}

func (this *Controller) Init() {}

func (this *Controller) SuperInit() {
	if !this.isInitialized {
		this.postFormParsed = false
		this.Data = make(map[interface{}]interface{})
		this.TplName = ""
		this.TplExt = "html"
		this.isInitialized = true
	}
}

func (this *Controller) Get() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Post() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Delete() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Put() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Head() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Patch() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Options() error {
	return errors.New("Not Allowed")
}

func (this *Controller) Render() error {
	return grend.HTML(this.W, http.StatusOK, this.TplName+"."+this.TplExt, this.Data)
}

func (this *Controller) ParseForm(maxMem int64) error {
	return this.R.ParseMultipartForm(maxMem)
}

func (this *Controller) PostFormValue(key string) string{
	return this.R.PostFormValue(key)
}

func (this *Controller) FormFile(key string) (multipart.File, *multipart.FileHeader, error){
	return this.R.FormFile(key)
}

//Initialize Global Server stuff
func Init(tplGlob string, debug bool, partialGlob string, funcs template.FuncMap) {
	grend = grender.New(grender.Options{
		Debug: debug,       // If true, templates will be recompiled before each render call
		TemplatesGlob: tplGlob,  // Glob to your template files
		PartialsGlob: partialGlob,   // Glob to your patials or global templates
		Funcs: funcs,         // Your template FuncMap
		Charset: "UTF-8",   // Charset to use for Content-Type header values
	})
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	grend.HTML(w, http.StatusForbidden, "forbidden.html", nil)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	grend.HTML(w, http.StatusNotFound, "notfound.html", nil)
}

type ListController struct {
	Controller
	ItemList interface{}
	ItemName string
	//Storer 		StoreInterface
}

func (this *ListController) Get() error {
	this.Data[this.ItemName] = this.ItemList
	return this.Render()
}


