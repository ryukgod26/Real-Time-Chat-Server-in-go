package main

import (
	"flag"
	"log"
	"net/http"
    "github.com/rivo/tview"
)

var addr = flag.String("addr",":8800","http Service addr")

func serveHome(w http.ResponseWriter, r *http.Request){
    log.Printf("[HTTP] %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	if r.URL.Path != "/"{
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet{
		http.Error(w, "Request Method is not GET", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w,r,"test.html")
}

func main(){
	flag.Parse()
	app := tview.NewApplication()
	root := buildMenu(app)
	
	if err := app.SetRoot(root, true).EnableMouse(true).Run(); err != nil{
		log.Fatal(err)
	}
	// hub := createHub()
	// go hub.run()
	// http.HandleFunc("/", serveHome)
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// err := http.ListenAndServe(*addr, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}
