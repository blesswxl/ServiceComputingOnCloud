package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/", testHandler(formatter)).Methods("GET")
}

func testHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, _ := readFile("data/data.json")
		req.ParseForm()
		// fmt.Println(req.Form)
		// fmt.Println(req.Form["url_long"])
		for name, _ := range req.Form {
			id, exist := data[name]
			if exist {
				fmt.Fprintf(w, id+"\n")
			} else {
				fmt.Fprintf(w, "This name does not exist!\n")
			}
		}
	}
}

func readFile(filename string) (map[string]string, error) {
	var data map[string]string
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return nil, err
	}

	return data, nil
}
