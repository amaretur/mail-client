package main

import (
	"github.com/amaretur/mail-client/app"
)

func main() {

	a := app.New()
	a.Init()
	a.Run()
}
