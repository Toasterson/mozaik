package main

import "github.com/toasterson/mozaik/models"

var tiles = map[int]Tile{
	0: {
		"Testing",
		sample_text,
		[]Tile{},
		[]string{"testPic"},
		[]string{},
		models.User{},
		TileNew,
		Asessment{},
	},
	1: {
		"Testing2",
		sample_text,
		[]Tile{},
		[]string{"testPic"},
		[]string{},
		models.User{},
		TileNew,
		Asessment{},
	},
}

var sample_text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
