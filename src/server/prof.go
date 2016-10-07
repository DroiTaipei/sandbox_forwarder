package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func ActivateProfile() {
	fmt.Println("Start profiling")
	go http.ListenAndServe(":10201", http.DefaultServeMux)
}
