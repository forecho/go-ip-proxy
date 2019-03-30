package server

import (
	"encoding/json"
	"fmt"
	"go-ip-proxy/config"
	"go-ip-proxy/storage"
	"net/http"
	"strings"
)

type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

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
	fmt.Printf(r.Method)
	if r.Method == "GET" {
		w.Header().Add("content-type", "application/json")
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
		w.Header().Add("content-type", "application/json")
		if s == nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		result := s.GetAll()
		if len(result) > 0 {
			//please use ',' as sep when is saved
			list := strings.Split(string(result),"http://")
			data := map[string]interface{}{
				"list": list,
			}
			JsonSuccess(w, data)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// deleteIp will delete the given ip. Return 200 if succeed.
// Sample usage: http://localhost:18090/delete?ip=0.0.0.0
func deleteIp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		values := r.URL.Query()
		if len(values["ip"]) > 1 {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if s.Delete(values["ip"][0]) {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func JsonSuccess(w http.ResponseWriter, data map[string]interface{}) {
	res := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	JsonReturn(w, res)
}

func JsonReturn(w http.ResponseWriter, res Response) {
	w.Header().Set("Content-Type", "application/json")
	str, _ := json.Marshal(res)
	w.Write(str)
}
