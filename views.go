package main

import (
	"github.com/toasterson/mozaik/router"
	"github.com/toasterson/mozaik/controller"
	"net/http"
	"github.com/gorilla/mux"
	"errors"
	"strconv"
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/mozaik/models"
	"github.com/toasterson/mozaik/logger"
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
type TileListController struct {
	controller.ListController
}

func (c *TileListController) Init() {
	c.SuperInit()
	c.TplName = "tile_list"
	c.ItemName = "tiles"
	c.ItemList = &tiles
	c.Routes = []router.Route{
		{http.MethodGet, "/tiles"},
	}
}

type TileDetailController struct {
	controller.Controller
}

func (c *TileDetailController) Init() {
	c.SuperInit()
	c.TplName = "tile_view"
	c.Routes = []router.Route{
		{http.MethodGet, "/tiles/{id}"},
	}
}

func (c *TileDetailController) Get() error  {
	vars := mux.Vars(c.R)
	id, err := strconv.Atoi(vars["id"])
	util.Must(err)
	//TODO Load from DB Here
	var ok bool
	c.Data["tile"], ok = tiles[id]
	if !ok {
		return errors.New("Does Not Exist yet Maybe create it?")
	}
	return c.Render()
}

type NewTileController struct {
	controller.Controller
}

func (c *NewTileController) Init() {
	c.SuperInit()
	c.TplName = "tile_form"
	c.EnableAuthProtection()
	c.Routes = []router.Route{
		{http.MethodGet, "/tiles/new"},
		{http.MethodPost, "/tiles/new"},
	}
}

func (c *NewTileController) Get() error {
	return c.Render()
}

func (c *NewTileController) Post() error {
	name := c.PostFormValue("name")
	p := Tile{
		name,
		[]byte(c.PostFormValue("text")),
		[]Tile{},
		[]string{},
		[]string{},
		models.User{},
		0,
		Asessment{},
	}
	//TODO Database Saving Here
	key := len(tiles)
	tiles[key] = p
	logger.Trace(tiles)
	http.Redirect(c.W, c.R, strconv.Itoa(key), http.StatusSeeOther)
	return nil
}