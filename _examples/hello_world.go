//go:build ignore
// +build ignore

package main

import (
	"log"
	"strconv"

	ui "github.com/deusnefum/termui/v3"
	"github.com/deusnefum/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	ui.StyleParserColorMap["240"] = ui.Color(240)
	p := widgets.NewParagraph()
	p.Title = "³³[³³](bg:240)³³He³llo³"
	p.TitleAlignment = ui.AlignCenter
	p.Text = "³Hello World!³"
	p.SetRect(10, 10, 35, 15)

	ui.Render(p)
	p.Text = strconv.Itoa(p.Max.X)
	ui.Render(p)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
	}
}
