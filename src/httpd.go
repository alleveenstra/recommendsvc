package main

import (
  	"recommendsvc"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log")

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage %s [input-file.json] [listen-address]", os.Args[0])
		os.Exit(1)
	}
	filebytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Unacle to read input file (%v)\n", err)
		return
	}
	log.Println("Loading places...")
	var places []recommendsvc.Place
	json_err := json.Unmarshal(filebytes, &places)
	if json_err != nil {
		log.Fatalf("JSON error: %v", json_err)
		return
	}
	log.Println("Starting webserver...")
	http.HandleFunc("/recommendation", recommendsvc.Build_recommendation_handler(places))
	http_err := http.ListenAndServe(os.Args[2], nil)
	if (http_err != nil) {
		log.Fatalf("Error creating server: %v", http_err)
		return
	}

}
