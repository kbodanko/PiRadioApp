package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
    zmq "github.com/pebbe/zmq4"
)

var client *zmq.Socket
var checkedStation string = "none"
var DEBUG bool = true

type StationsName struct {
    Name string
}

type StationsPageData struct {
    Stations []StationsName
}

type PlayingDetails struct {
    Station string
    Title string
    Url string
}

type ViewPageData struct {
    StationsData StationsPageData
    PlayingData PlayingDetails
}

func getStations() []StationsName {
    var stations []StationsName
    strStations := sendRequest("station get all")
    for i := range strStations {
        stations = append(stations, StationsName{strStations[i]})
    }
    return stations
}

/* 
    HANDLERS: api 
*/

func viewHandler(w http.ResponseWriter, r *http.Request){
    debugMsg("main page")
    tmpl, err := template.ParseFiles("../server/server.html")
    if err != nil {
        debugMsg("Error occure: " + err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    current := sendRequest("station get current")
    details := ViewPageData {
        StationsData: StationsPageData{
            Stations: getStations(),
        },
        PlayingData: PlayingDetails {
            Station: current[0],
            Title: current[1],
            Url: current[2],
        },
    }

    tmpl.Execute(w, details)
}

func playHandler(w http.ResponseWriter, r *http.Request){
    sendRequest("engine state set play")
    debugMsg("play ")
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

func stopHandler(w http.ResponseWriter, r *http.Request){
    sendRequest("engine state set stop")
    debugMsg("stop")
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

func setHandler(w http.ResponseWriter, r *http.Request){
    body := r.FormValue("Stations")
    checkedStation = body
    debugMsg("checked " + checkedStation)
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

func submitHandler(w http.ResponseWriter, r *http.Request){
    sendRequest("station set " + checkedStation)
    sendRequest("engine state set play")
    debugMsg("submit " + checkedStation)
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

func deleteHandler(w http.ResponseWriter, r *http.Request){
    sendRequest("station delete " + checkedStation)
    debugMsg("delete " + checkedStation)
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

func addHandler(w http.ResponseWriter, r *http.Request){
    name := r.FormValue("stationName")
    uri := r.FormValue("stationUri")
    sendRequest("station new " + name + " " + uri)
    debugMsg("add new " + name + " " + uri)
    http.Redirect(w, r, "/radio/", http.StatusFound)
} 

/*
    TCP connection with engine
*/

func sendRequest(request string) []string {
    /*
        send request; recive, decode and return reply
    */
    (*client).SendMessage(request, 0)
    debugMsg("send request: " + request)
    msg, err := (*client).RecvMessage(0)
    if err != nil {
        debugMsg("send error: " + err)
    }
    reply := msg[0]
    debugMsg("recive replay: " + reply)
    var words []string
	word := ""
	flag := true
	iter := -1
	for i := range reply {
		if reply[i] == ' ' {
			flag = true
		} else if flag {
			word = ""
			words = append(words, word)
			iter++
			words[iter] = words[iter] + string(reply[i])
			flag = false
		} else {
			words[iter] = words[iter] + string(reply[i])
		}
	}
    return words 
}

func debugMsg(message string){
    if DEBUG {
        fmt.Println(message)
    }
}

func main() {
    socket, errSocket := zmq.NewSocket(zmq.PAIR)
    client = socket
    engineUri = "tcp://localhost:5555"
    serveUri = ":8080"
    (*client).Connect(engineUri)
    if errServe != nil {
        debugMsg("error socket: " + errSocket)
    } else {
        debugMsg("connected to server: " + engineUri)
    }
    
    http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../server/"))))
    http.HandleFunc("/radio/play", playHandler)
    http.HandleFunc("/radio/stop", stopHandler)
    http.HandleFunc("/radio/set", setHandler)
    http.HandleFunc("/radio/submit", submitHandler)
    http.HandleFunc("/radio/delete", deleteHandler)
    http.HandleFunc("/radio/add", addHandler)
    http.HandleFunc("/radio", viewHandler)

    debugMsg("serve uri: " + serveUri)
    errServe := http.ListenAndServe(serveUri, nil))
    if errServe != nil {
        debugMsg("error serve: " + err) 
    }
}
