package game

type Element struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MapNode struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	X      int       `json:"x"`
	Y      int       `json:"y"`
	Worlds []Element `json:"worlds"`
}

type WorldNodes map[string]MapNode

type World struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Points WorldNodes `json:"points"`
}

var Universe = map[string]World{
	"main_city": {
		ID:   "main_city",
		Name: "Благословенный край",
		Points: WorldNodes{
			"start_glade": {
				ID:   "start_glade",
				Name: "Дом родной",
				X:    363, Y: 185,
				Worlds: []Element{},
			},
			"tample": {
				ID:   "tample",
				Name: "Храм",
				X:    194, Y: 180,
				Worlds: []Element{},
			},
			"tavern": {
				ID:   "tavern",
				Name: "Таверна",
				X:    195, Y: 269,
				Worlds: []Element{},
			},
			"arena": {
				ID:   "arena",
				Name: "Арена",
				X:    243, Y: 128,
				Worlds: []Element{},
			},
			"shop": {
				ID:   "shop",
				Name: "Лавка",
				X:    406, Y: 272,
				Worlds: []Element{},
			},
			"bestiary": {
				ID:   "bestiary",
				Name: "Тренировочная",
				X:    193, Y: 99,
				Worlds: []Element{},
			},
			"academy": {
				ID:   "academy",
				Name: "Академия",
				X:    284, Y: 241,
				Worlds: []Element{},
			},
			"portal": {
				ID:   "portal",
				Name: "Портал",
				X:    449, Y: 125,
				Worlds: []Element{{ID: "totlhaim", Name: "Тотльхайм"}},
			},
			"armory": {
				ID:   "armory",
				Name: "Оружейная",
				X:    323, Y: 101,
				Worlds: []Element{},
			},
		},
	},
	"totlhaim": { // Второй мир Тотльхайм
		ID:   "totlhaim",
		Name: "Тотльхейм",
		Points: WorldNodes{
			"portal": {
				ID:   "portal",
				Name: "Портал",
				X:    449, Y: 125,
				Worlds: []Element{{ID: "main_city", Name: "Благословенный край"}},
			},
		},
	},
}
