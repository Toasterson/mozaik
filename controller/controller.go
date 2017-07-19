package controller

import (
	"net/http"
	"errors"
	"github.com/toasterson/mozaik/router"
	"mime/multipart"
	"gopkg.in/authboss.v1"
	"reflect"
	"github.com/toasterson/mozaik/logger"
	"fmt"
	"github.com/justinas/nosurf"
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
	Auth *authboss.Authboss
	isProtected bool
	RedirectUrl string
}

func (this *Controller) IsAuthProtected() bool{
	return this.isProtected
}

func (this *Controller) EnableAuthProtection(){
	this.isProtected = true
}

func (this *Controller) GetAuth() *authboss.Authboss{
	return this.Auth
}

func (this *Controller) GetRoutes()	[]router.Route{
	return this.Routes
}

func (this *Controller) Prepare(w http.ResponseWriter, r *http.Request){
	this.W = w
	this.R = r
	//TODO add Route to this Data Value
	for key, val := range layoutData(w,r){
		this.Data[key] = val
	}
	this.Data["xsrfName"] = "csrf_token"
	this.Data["xsrfToken"] = nosurf.Token(r)
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
		this.isProtected = false
		this.Auth = AuthBoss
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
	logger.Trace(fmt.Sprintf("Rendering %s with data %s", this.TplName, this.Data))
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

func (this *Controller) AddData(dat map[interface{}]interface{}){
	for key, val := range dat {
		this.Data[key] = val
	}
}

func (this *Controller) AddHTMLData(dat map[string]interface{}){
	for key, val := range dat {
		this.Data[key] = val
	}
}

func getName(controller interface{}) string{
	if t := reflect.TypeOf(controller); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}