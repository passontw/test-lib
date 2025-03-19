package console

import (
	"github.com/common-nighthawk/go-figure"
)

func Print(text string) {
	myFigure := figure.NewFigure(text, "", true)
	myFigure.Print()
}
