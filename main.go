package main

import (
	"html/template"
	"log"
	"net/http"
)

type state bool

var (
	ADDRESS       = ":8080"
	PIN           = "7"
	ON      state = true
	OFF     state = false
	STATE         = OFF
)

func init() {
	setPinState(OFF)
}

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/on", on)
	mux.HandleFunc("/off", off)

	log.Printf("listening on %s\n", ADDRESS)
	http.ListenAndServe(ADDRESS, mux)
}

type Params struct {
	State string
	URL   string
}

func NewParams(value state) *Params {
	var (
		state string
		url   string
	)
	switch value {
	case ON:
		state = "ON"
		url = "on"
	case OFF:
		state = "OFF"
		url = "off"
	}
	return &Params{
		State: state,
		URL:   url,
	}
}

func on(w http.ResponseWriter, r *http.Request) {
	setPinState(ON)
	t, _ := template.ParseFiles("switch.html")
	t.Execute(w, NewParams(!STATE))
}

func off(w http.ResponseWriter, r *http.Request) {
	setPinState(OFF)
	t, _ := template.ParseFiles("switch.html")
	t.Execute(w, NewParams(!STATE))
}

func setPinState(s state) {
	STATE = s
	stateValue := ""
	switch s {
	case ON:
		stateValue = "high"
	case OFF:
		stateValue = "low"
	}
	log.Printf("setting state to %s\n", stateValue)
}
