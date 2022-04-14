// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"
	"math"

	. "github.com/deusnefum/termui/v3"
)

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []float64
	Title      string
	TitleStyle Style
	LineColor  Color
	MaxVal     float64
	MaxHeight  int // TODO
}

// SparklineGroup is a renderable widget which groups together the given sparklines.
type SparklineGroup struct {
	Block
	Sparklines []*Sparkline
}

// NewSparkline returns a unrenderable single sparkline that needs to be added to a SparklineGroup
func NewSparkline() *Sparkline {
	return &Sparkline{
		TitleStyle: Theme.Sparkline.Title,
		LineColor:  Theme.Sparkline.Line,
	}
}

func NewSparklineGroup(sls ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{
		Block:      *NewBlock(),
		Sparklines: sls,
	}
}

func (self *SparklineGroup) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	sparklineHeight := self.Inner.Dy() / len(self.Sparklines)

	for i, sl := range self.Sparklines {
		heightOffset := (sparklineHeight * (i + 1))
		barHeight := sparklineHeight
		if i == len(self.Sparklines)-1 {
			heightOffset = self.Inner.Dy()
			barHeight = self.Inner.Dy() - (sparklineHeight * i) - 1
		}
		if sl.Title != "" && i != len(self.Sparklines)-1 {
			barHeight--
		}

		maxVal := sl.MaxVal
		if maxVal == 0 {
			maxVal, _ = GetMaxFloat64FromSlice(sl.Data)
		}

		// draw line
		for j := 0; j < len(sl.Data) && j < self.Inner.Dx(); j++ {
			data := sl.Data[j]
			if math.IsNaN(data) {
				data = 0
			}
			height := int((data / maxVal) * float64(barHeight))
			sparkChar := BARS[len(BARS)-1]
			for k := 0; k < height+1; k++ {
				buf.SetCell(
					NewCell(sparkChar, NewStyle(sl.LineColor)),
					image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-k),
				)
			}
			heightBlocksCnt := len(BARS) - 1
			lastHeight := int((data/maxVal)*float64(barHeight*heightBlocksCnt)) % heightBlocksCnt
			// prevent gaps from showing if at bottom of sparkline
			if lastHeight == 0 && height == 0 {
				lastHeight = 1
			}
			buf.SetCell(
				NewCell(BARS[lastHeight], NewStyle(sl.LineColor)),
				image.Pt(j+self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-height),
			)
		}

		if sl.Title != "" {
			// draw title
			buf.SetString(
				TrimString(sl.Title, self.Inner.Dx()),
				sl.TitleStyle,
				image.Pt(self.Inner.Min.X, self.Inner.Min.Y-1+heightOffset-barHeight),
			)
		}
	}
}
