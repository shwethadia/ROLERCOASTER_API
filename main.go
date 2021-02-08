package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

//Coaster ...
type Coaster struct {

	Name string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	ID string	`json:"id"`
	Height int	`json:"height"`
	InPark string	`json:"inpark"`
}

type coasterHandlers struct {

	sync.Mutex
	store map[string]Coaster
}

type adminPortal struct {

	password string
}

func newAdminPortal() *adminPortal{

	password := os.Getenv("ADMIN_PASSWORD")
	if password == ""{
		panic("Required env var ADMIN PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter,r *http.Request){

	user,pass,ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password{

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized"))
		return 
	}

	w.Write([]byte("<html><h1>Super Secret Admin Portal</h1></html>"))
}

func main(){


	admin := newAdminPortal()
	coasterHandlers := newCoasterHandlers()
	http.HandleFunc("/coasters",coasterHandlers.coasters)
	http.HandleFunc("/coasters/",coasterHandlers.getcoaster)
	http.HandleFunc("/admin",admin.handler)
	err := http.ListenAndServe(":8080",nil)
	if err!=nil{
		panic(err)
	}
}

func (c *coasterHandlers) coasters(w http.ResponseWriter,r *http.Request){


	switch r.Method{
	case "GET":
		c.get(w,r)
		return
	case "POST":
		c.post(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return 
	}
}

func (c *coasterHandlers) post(w http.ResponseWriter,r *http.Request){

bodyBytes,err := ioutil.ReadAll(r.Body)
defer r.Body.Close()
if err != nil {

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	return 
} 

ct := r.Header.Get("content-type")
if ct != "application/json"{

	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Write([]byte(fmt.Sprintf("Need content-type 'application/json',but got '%s' ",ct)))
	return 
}
var coaster Coaster
err = json.Unmarshal([]byte(bodyBytes),&coaster)
if err != nil {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
	return 
}

coaster.ID = fmt.Sprintf("%d",time.Now().UnixNano())
c.Lock()
c.store[coaster.ID] = coaster
defer c.Unlock()

}

func (c *coasterHandlers) get(w http.ResponseWriter,r *http.Request){


	coasters := make([]Coaster,len(c.store))

	c.Lock()
	i := 0
	for _,coaster :=  range c.store{

		coasters[i] = coaster
		i++
	}
	c.Unlock()

	jsonBytes, err := json.MarshalIndent(coasters,"","  ")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} 
	w.Header().Add("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

//getRandomCoaster ...
func (c *coasterHandlers) getRandomCoaster(w http.ResponseWriter,r *http.Request){

	ids := make([]string,len(c.store))
	c.Lock()
	i := 0
	for id := range c.store{

		ids[i] = id
		i++
	}
	defer c.Unlock()

	var target string
	if len(ids) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}else if len(ids) == 1 {
		target =ids[0]
	}else{
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}
	w.Header().Add("location",fmt.Sprintf("/coasters/%s",target))
	w.WriteHeader(http.StatusFound)

}
func (c *coasterHandlers) getcoaster(w http.ResponseWriter,r *http.Request){


	parts := strings.Split(r.URL.String(),"/")
	if len(parts) !=3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {

		c.getRandomCoaster(w,r)
		return 

	}
	c.Lock()
	coaster,ok := c.store[parts[2]]
	c.Unlock()
	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(coaster,"","  ")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} 
	w.Header().Add("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func newCoasterHandlers() *coasterHandlers{

	return &coasterHandlers{
		store: map[string]Coaster{
			"id1":Coaster{
				Name : "Fury 325",
				Height:99,
				ID:"id1",
				InPark: "Carowinds",
				Manufacturer: "B-M",
			},
		},
	}
}