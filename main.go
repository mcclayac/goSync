package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"poetry"
	"strconv"
)

/*
Anthonys-MacBook-Pro:go mcclayac$ godoc fmt Fprintf | more
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
Fprintf formats according to a format specifier and writes to w. It
returns the number of bytes written and any write error encountered.
*/

type poemWithTitle struct {
	Title     string
	Body      poetry.Poem
	WordCount string
	TheCount  int
}

var cache map[string]poetry.Poem

type config struct {
	Route       string
	BindAddress string   `json:"addr"`
	ValidPoems  []string `json:"valid"`
	//{"doggie","cat","letterA.txt"}

}

var c config

func poemHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	poemName := r.Form["name"][0]
	//fileName := "doggie"

	p, ok := cache[poemName]
	//p, err := poetry.LoadPoem(poemName)

	if !ok {
		http.Error(w, "File Not Found", http.StatusInternalServerError)
		fmt.Printf("An Error occured reading file %s \n", poemName)
		//os.Exit(-1)
		return
	}

	log.Printf("user Requested poem %s \n", poemName)

	p.ShufflePoem()
	pwt2 := poemWithTitle{poemName, p,
		strconv.FormatInt(int64(p.NumWords()), 10),
		p.NumThe()}
	enc2 := json.NewEncoder(w)

	enc2.Encode(pwt2)
}

func main() {
	log.SetFlags(log.Lmicroseconds)

	f, err := os.Open("config")
	if err != nil {
		log.Fatalf("Connot open the config file\n")
		//os.Exit(-1)
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	err = dec.Decode(&c)
	if err != nil {
		log.Fatalf("Cannot De-code the JSON config file")
		//os.Exit(-1)
	}
	p := poetry.Poem{}

	// Must make the map before you use it !!!!!
	cache = make(map[string]poetry.Poem)

	// Sync - read in all poems
	for _, name := range c.ValidPoems {
		p, err = poetry.LoadPoem(name)

		cache[name] = p
		if err != nil {
			log.Fatalf("Error Loading Poems") // Log Fatal message
			//os.Exit(-1)
		}
	}

	fmt.Printf("%v\n\n", c)

	http.HandleFunc(c.Route, poemHandler)
	http.ListenAndServe(c.BindAddress, nil)

}

//  curl http://127.0.0.1:8088/get\?name=cat | json_pp
