package cvi_listener

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)

const (
	API_PREFIX = "/api/v1"
)

type ApiServer struct {
	cvi3_manager *CVI3Manager
}

func (as *ApiServer) StartService(manager *CVI3Manager) {
	as.cvi3_manager = manager

	m := mux.NewRouter()

	m.HandleFunc(API_PREFIX + "/api-doc", APIDoc).Methods("GET")
	m.HandleFunc(API_PREFIX + "/psets", PSets).Methods("PUT")
	m.HandleFunc(API_PREFIX + "/status", Status).Methods("GET")
	m.HandleFunc(API_PREFIX + "/results", Results).Methods("GET")

	err := http.ListenAndServe(":8080", m) //设置监听的端口
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func APIDoc(w http.ResponseWriter, r *http.Request) {

}

func PSets(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	fmt.Printf("%s\n", v)
}

func Status(w http.ResponseWriter, r *http.Request) {

}

func Results(w http.ResponseWriter, r *http.Request) {

}