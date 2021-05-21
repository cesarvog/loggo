package main

import (
	"os"
    "fmt"
    "log"
	"io/ioutil"
    "net/http"
    "github.com/gorilla/mux"
	"time"
	"strconv"
)


var fmngr *FileManager
var w chan string

func handleRequests(p string) {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/tail/{qtd}", tail)
    myRouter.HandleFunc("/info", info)
    myRouter.HandleFunc("/warn", warn)
    myRouter.HandleFunc("/error", erro)

    log.Fatal(http.ListenAndServe(":"+p, myRouter))
}
		
func main() {
	fmt.Println("loggo - Log Service")
	port := os.Args[1]
	path := os.Args[2]

	if port == "" {
		fmt.Println("what port?")
		os.Exit(1)
	}

	if path == "" {
		fmt.Println("what file path to log in?")
	}

	fmngr = NewFileManager(path)
	w = fmngr.WriteChan()

	go func() { 
		fmngr.Run() 
	}()

	handleRequests(port)
}

func tail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qtd, err := strconv.Atoi(vars["qtd"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t, err := fmngr.Tail(qtd)
	if err != nil {
		t = err.Error()
	}

	fmt.Fprintf(w, t)
}

func info(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		logg("INFO", string(body))
	} else {
		fmt.Println(err.Error())
	}

}

func warn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		logg("WARN", string(body))
	} else {
		fmt.Println(err.Error())
	}

}

func erro(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		logg("ERROR", string(body))
	} else {
		fmt.Println(err.Error())
	}
}

func logg(verb string, t string) {
	dt := time.Now().Format("01-02-2006 15:04:05")
	w <- fmt.Sprintf("%s %s %s\n", dt, verb, t) 
}
