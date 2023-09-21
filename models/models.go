package models

type board struct {
	Title  string  `json:"title"`
	Fields []field `json:"fields"`
}

type field struct {
	Name     string   `json:"name"`
	Low      int      `json:"low"`
	High     int      `json:"high"`
	Tags     []string `json:"tags"`
	ImageUrl string   `json:"image_url"`
	Points   int      `json:"points"`
	Id       int      `json:"id"`
	Games    []string `json:"games"`
	Rating   int      `json:"r ating"`
}
