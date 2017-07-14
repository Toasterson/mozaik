package main

import (
	"github.com/toasterson/mozaik/router"
	"github.com/toasterson/mozaik/controller"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"errors"
)

type MainController struct {
	controller.Controller
}

func (c *MainController) Get() error {
	return c.Render()
}

func (c *MainController) Init() {
	c.SuperInit()
	c.TplName = "index"
	c.Routes = []router.Route{
		{http.MethodGet, "/"},
	}
}

//TODO See if we can make one ListController with different Parameters instead of one per Object
type PageDetailController struct {
	controller.Controller
}

func (c *PageDetailController) Init() {
	c.SuperInit()
	c.TplName = "page_view"
	c.Routes = []router.Route{
		{http.MethodGet, "/pages/{title}"},
	}
}

func (c *PageDetailController) Get() error  {
	vars := mux.Vars(c.R)
	title := vars["title"]
	//TODO Load from DB Here
	var ok bool
	c.Data["page"], ok = pages[title]
	if !ok {
		return errors.New("Does Not Exist yet Maybe create it?")
	}
	return c.Render()
}

type NewPageController struct {
	controller.Controller
}

func (c *NewPageController) Init() {
	c.SuperInit()
	c.TplName = "page_form"
	c.Routes = []router.Route{
		{http.MethodGet, "/pages/new"},
		{http.MethodPost, "/pages/new"},
	}
}

func (c *NewPageController) Get() error {
	return c.Render()
}

func (c *NewPageController) Post() error{
	title := c.PostFormValue("title")
	p := Page{title, []byte(c.PostFormValue("body"))}
	//TODO Database Saving Here
	pages[title] = p
	fmt.Println(pages)
	http.Redirect(c.W, c.R, title, http.StatusSeeOther)
	return nil
}
