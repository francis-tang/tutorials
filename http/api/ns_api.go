package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

type NameServiceApi struct {
	http.Handler
}

func NewNameServiceApi() *NameServiceApi{
	r:=mux.NewRouter()
	nsa := &NameServiceApi{
		Handler:r,
	}
	r.StrictSlash(true)
	r.Methods("GET").Path("/{key}").HandlerFunc(nsa.handleGet)
	r.Methods("POST").Path("/{key}").HandlerFunc(nsa.handleSet)

	return nsa
}

func (nsa *NameServiceApi) handleGet(w http.ResponseWriter,r *http.Request){
	key:=mux.Vars(r)["key"]
	if key == ""{
		respond(w,http.StatusBadRequest,apiResponse{Error:"no key"})
		return
	}

	respond(w,http.StatusOK,apiResponse{
		Key:key,
		Value:"ok,"+key,
	})
}

func (nsa *NameServiceApi) handleSet(w http.ResponseWriter,r *http.Request){
	key := mux.Vars(r)["key"]
	if key == ""{
		respond(w,http.StatusBadRequest,apiResponse{Error:"no key"})
		return
	}

	respond(w,http.StatusOK,apiResponse{
		Key:key,
		Value:"good,"+key,
	})
}

type apiResponse struct {
	Key string     `json:"key,omitempty"`
	Value string   `json:"value,omitempty"`
	Error string   `json:"error,omitempty"`
	Info string    `json:"info,omitempty"`
	Log string     `json:"log,omitempty"`
	//Key   string `json:"key,omitempty"`
	//Value string `json:"value,omitempty"`
	//Error string `json:"error,omitempty"`
	//Info  string `json:"info,omitempty"`
	//Log   string `json:"log,omitempty"`
}

func respond(w http.ResponseWriter,code int,response apiResponse){
	w.WriteHeader(code)
	buf,_ := json.MarshalIndent(response,"","    ")
	w.Write(buf)
}
