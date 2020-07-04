package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)


type Event struct {
	ResizeFrom		Dimension
	ResizeTo			Dimension
	WebsiteUrl 		string		`json:"siteUrl"`
	Pasted				bool
	Time 					int
	SessionId 		string
	FormId 				string
	EventType 		string
}

type server struct{}


var clientSessions map[string]*Data


func main() {
	s := &server{}
	http.Handle("/api", s)
	clientSessions = make(map[string]*Data)
	http.HandleFunc("/api/events", eventApi)
	
	//helper endpoint for testing
	http.HandleFunc("/api/sessions/{id}", sessionApi)

	log.Fatal(http.ListenAndServe(":8080", nil))

}


func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Del("Content-Type")
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		w.WriteHeader(http.StatusOK)

    w.Write([]byte(`{"message": "Use POST request on /api/event API endpoint to send Event data \n Use GET request on /api/session/{session-id} to get full session stored by now"`))
}


func eventApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

		var e Event
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&e)
		

		if err != nil {
			responseError(w, http.StatusBadRequest, "Error parsing JSON")
		}
		if e.SessionId == "" {
			responseError(w, http.StatusBadRequest, "Invalid session id")
		}

		s := clientSessions[e.SessionId]
		if s == nil {
			s = NewData(e)
		}

		err = s.updateSession(e)
		if err != nil {
			responseError(w, http.StatusBadRequest, err.Error())
		} 
		clientSessions[e.SessionId] = s
		if s.isCompleted() {
			fmt.Println("Form submitted, struct completed")

			delete(clientSessions, s.SessionId)
		}
		s.Print()

		w.WriteHeader(http.StatusCreated)
    w.Write([]byte(""))
	} else {
		w.WriteHeader(http.StatusOK)
    w.Write([]byte(""))
	}
}


func sessionApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		sid := strings.TrimPrefix(r.URL.Path, "/api/")
    if sid == "" {
			responseError(w, http.StatusBadRequest, "Invalid session id")
		}
		s := clientSessions[sid]
		if s == nil {
			responseError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*s)

	} else if r.Method == "OPTIONS" || r.Method == "HEAD" {
		//for CORS to work properly  OPTIONS has to return OK as well along with with headers
		w.WriteHeader(http.StatusOK)
    w.Write([]byte(""))
	} else {
		responseError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}


func responseError(w http.ResponseWriter, code int, m string) {
	w.WriteHeader(code)
	http.Error(w, m, code)

}