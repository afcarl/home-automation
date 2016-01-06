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
	setPinState(OFF, nil)
}

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)
	mux.HandleFunc("/on", on)
	mux.HandleFunc("/off", off)
	mux.HandleFunc("/home", home)

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

func home(w http.ResponseWriter, r *http.Request) {
	r.Header["X-Home-Automation"] = []string{"ignore"}
	switch STATE {
	case ON:
		on(w, r)
	case OFF:
		off(w, r)
	}
}

func on(w http.ResponseWriter, r *http.Request) {
	setPinState(ON, r.Header)
	t, _ := template.ParseFiles("switch.html")
	t.Execute(w, NewParams(!STATE))
}

func off(w http.ResponseWriter, r *http.Request) {
	setPinState(OFF, r.Header)
	t, _ := template.ParseFiles("switch.html")
	t.Execute(w, NewParams(!STATE))
}

func setPinState(s state, header http.Header) {
	stateValue := ""
	switch s {
	case ON:
		stateValue = "high"
	case OFF:
		stateValue = "low"
	}

	val, ok := header["X-Home-Automation"]
	if ok && val[0] == "ignore" {
		log.Printf("ignore pin set %s\n", stateValue)
		return
	}

	STATE = s
	log.Printf("setting state to %s\n", stateValue)
}
