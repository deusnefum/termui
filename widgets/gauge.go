// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"

	. "github.com/deusnefum/termui/v3"
)

type Gauge struct {
	Block
	Percent         float64
	BarColor        Color
	Label           string
	LabelStyle      Style
	LabelOnBarStyle Style
}

func NewGauge() *Gauge {
	return &Gauge{
		Block:           *NewBlock(),
		BarColor:        Theme.Gauge.Bar,
		LabelOnBarStyle: NewStyle(Theme.Gauge.Bar, ColorClear, ModifierReverse),
		LabelStyle:      Theme.Gauge.Label,
	}
}

func (self *Gauge) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	label := self.Label
	if label == "" {
		label = fmt.Sprintf("%.0f%%", self.Percent)
	}

	// plot bar
	shadeCnt := len(SHADED_BLOCKS) - 1
	barWidth := int(self.Percent * float64(self.Inner.Dx()))
	lastBarWidth := int(self.Percent*float64(self.Inner.Dx()*shadeCnt)) % shadeCnt
	if lastBarWidth < 0 {
		lastBarWidth = 0
	}
	if barWidth > 0 {
		buf.Fill(
			NewCell(' ', NewStyle(ColorClear, self.BarColor)),
			image.Rect(self.Inner.Min.X, self.Inner.Min.Y, self.Inner.Min.X+barWidth, self.Inner.Max.Y),
		)
	}
	if self.Inner.Min.X+barWidth+1 < self.Inner.Max.X {
		buf.Fill(
			NewCell(SHADED_BLOCKS[lastBarWidth], NewStyle(self.BarColor, ColorClear)),
			image.Rect(self.Inner.Min.X+barWidth, self.Inner.Min.Y, self.Inner.Min.X+barWidth+1, self.Inner.Max.Y),
		)
	}

	// plot label
	labelXCoordinate := self.Inner.Min.X + (self.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := self.Inner.Min.Y + ((self.Inner.Dy() - 1) / 2)
	if labelYCoordinate < self.Inner.Max.Y {
		for i, char := range label {
			if labelXCoordinate+i+1 <= self.Inner.Min.X+barWidth {
				buf.SetCell(NewCell(char, self.LabelOnBarStyle), image.Pt(labelXCoordinate+i, labelYCoordinate))
				continue
			}
			buf.SetCell(NewCell(char, self.LabelStyle), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
