package utils

import (
	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// StdoutBanner prints a banner with the given text to stdout.
func StdoutBanner(pro string) {
	color.Green("%s \n\n", RenderBanner(pro))
}

// RenderBanner renders a banner with the given text.
func RenderBanner(pro string) string {
	banner := figure.NewFigure(pro, "", true)
	out := banner.String()
	return out
}
