package main

import (
	"AI4All/window"
)

// System Main Function
// to compile: CGO_LDFLAGS="-lopenblas" LIBRARY_PATH=$PWD C_INCLUDE_PATH=$PWD go run -tags openblas . -t 14

func main() {
	win := window.NewWindow("AI4All")
	win.Win.ShowAndRun()
}
