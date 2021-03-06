package server

import (
	"encoding/json"
	"fmt"
	"go-ip-proxy/config"
	"go-ip-proxy/result"
	"go-ip-proxy/storage"
	"net/http"
)

var s storage.Storage

// NewServer will start a new server.
func NewServer(storage storage.Storage) {
	if storage != nil {
		s = storage
	}

	http.HandleFunc("/get", getIp)
	http.HandleFunc("/get-all", getAll)
	http.HandleFunc("/delete", deleteIp)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config().Http.Port), nil)
	if err != nil {
		panic(err)
	}
}

// getIp will get a random Ip.
// Sample usage: http://localhost:18090/get
func getIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		if s == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		result := s.Get()
		if len(result) > 0 {
			w.Write([]byte(result))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getAll will get all Ip.
// Sample usage: http://localhost:18090/get-all
func getAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		if s == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		res := s.GetAll()
		if len(res) > 0 {
			newResult := result.Result{
				Data:  res,
				Count: len(res),
			}
			str, _ := json.Marshal(newResult)
			w.Write(str)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// deleteIp will delete the given ip. Return 200 if succeed.
// Sample usage: http://localhost:18090/delete?proxy=http://101.37.20.241:443
func deleteIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		values := r.URL.Query()
		if len(values["proxy"]) > 1 {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if s.Delete(values["proxy"][0]) {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
