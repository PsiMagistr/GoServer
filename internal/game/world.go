package game

type MapNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

/*var WorldNodes = map[string]MapNode{
	"start_glade": {
		ID:   "start_glade",
		Name: "Дом родной",
		X:    363, Y: 185,
	},
	"tample": {
		ID:   "tample",
		Name: "Храм",
		X:    194, Y: 180,
	},
	"tavern": {
		ID:   "tavern",
		Name: "Таверна",
		X:    195, Y: 269,
	},
	"arena": {
		ID:   "arena",
		Name: "Арена",
		X:    243, Y: 128,
	},
	"shop": {
		ID:   "shop",
		Name: "Лавка",
		X:    406, Y: 272,
	},
}*/

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
			},
			"tample": {
				ID:   "tample",
				Name: "Храм",
				X:    194, Y: 180,
			},
			"tavern": {
				ID:   "tavern",
				Name: "Таверна",
				X:    195, Y: 269,
			},
			"arena": {
				ID:   "arena",
				Name: "Арена",
				X:    243, Y: 128,
			},
			"shop": {
				ID:   "shop",
				Name: "Лавка",
				X:    406, Y: 272,
			},
			"bestiary": {
				ID:   "bestiary",
				Name: "Тренировочная",
				X:    193, Y: 99,
			},
			"academy": {
				ID:   "academy",
				Name: "Академия",
				X:    284, Y: 241,
			},
			"portal": {
				ID:   "portal",
				Name: "Портал",
				X:    449, Y: 125,
			},
			"armory": {
				ID:   "armory",
				Name: "Оружейная",
				X:    323, Y: 101,
			},
		},
	},
}
