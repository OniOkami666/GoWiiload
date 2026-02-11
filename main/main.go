package main

import (
	"flag"
	"fmt"

	"github.com/OniOkami666/GoWiiload/gowiiload"
)

func main() {
	file := flag.String("file", "", "Load a dol file")
	flag.Parse()

	if *file == "" {
		fmt.Println("Please enter a valid dol file!")
		return
	}
	err := gowiiload.WiiloadConnect("", *file)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		print("File sent")
	}
}
