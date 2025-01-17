// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"
	"image"
	"math"

	. "github.com/sparques/termui/v3"
)

// Plot has two modes: line(default) and scatter.
// Plot also has two marker types: braille(default) and dot.
// A single braille character is a 2x4 grid of dots, so using braille
// gives 2x X resolution and 4x Y resolution over dot mode.
type Plot struct {
	Block

	Data       [][]float64
	DataLabels []string
	MaxVal     float64
	MinVal     float64

	LineColors  []Color
	AxesColor   Color
	LabelStyle  Style
	ShowAxes    bool
	YAxisFormat string

	Marker          PlotMarker
	DotMarkerRune   rune
	PlotType        PlotType
	HorizontalScale float64
	DrawDirection   DrawDirection // TODO

	yAxisLabelsWidth int
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
	axisFormat        = "%.2f"
)

type PlotType uint

const (
	LineChart PlotType = iota
	ScatterPlot
)

type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

type DrawDirection uint

const (
	DrawLeft DrawDirection = iota
	DrawRight
)

func NewPlot() *Plot {
	return &Plot{
		Block:            *NewBlock(),
		LineColors:       Theme.Plot.Lines,
		AxesColor:        Theme.Plot.Axes,
		Marker:           MarkerBraille,
		DotMarkerRune:    DOT,
		Data:             [][]float64{},
		HorizontalScale:  1,
		DrawDirection:    DrawRight,
		ShowAxes:         true,
		LabelStyle:       NewStyle(ColorClear),
		PlotType:         LineChart,
		YAxisFormat:      axisFormat,
		yAxisLabelsWidth: yAxisLabelsWidth,
	}
}

func (self *Plot) renderBraille(buf *Buffer, drawArea image.Rectangle, maxVal, minVal float64) {
	canvas := NewCanvas()
	canvas.Rectangle = drawArea

	// don't do anything if there's no data
	if self.Data == nil || len(self.Data) == 0 {
		return
	}

	switch self.PlotType {
	case ScatterPlot:
		for i, line := range self.Data {
			for j, val := range line {
				height := int(((val - minVal) / (maxVal - minVal)) * float64(drawArea.Dy()-1))
				canvas.SetPoint(
					image.Pt(
						(drawArea.Min.X+(int(math.Round(float64(j)*self.HorizontalScale))))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					SelectColor(self.LineColors, i),
				)
			}
		}
	case LineChart:
		for i, line := range self.Data {
			// skip plotting a line with fewer than 2 data points
			if len(line) < 2 {
				continue
			}
			previousHeight := int(((line[0] - minVal) / (maxVal - minVal)) * float64(drawArea.Dy()-1))
			for j, val := range line[1:] {
				height := int(((val - minVal) / (maxVal - minVal)) * float64(drawArea.Dy()-1))
				canvas.SetLine(
					image.Pt(
						(drawArea.Min.X+(int(math.Round(float64(j)*self.HorizontalScale))))*2,
						(drawArea.Max.Y-previousHeight-1)*4,
					),
					image.Pt(
						(drawArea.Min.X+(int(math.Round(float64(j+1)*self.HorizontalScale))))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					SelectColor(self.LineColors, i),
				)
				previousHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

func (self *Plot) renderDot(buf *Buffer, drawArea image.Rectangle, maxVal, minVal float64) {
	switch self.PlotType {
	case ScatterPlot:
		for i, line := range self.Data {
			for j, val := range line {
				height := int(((val - minVal) / (maxVal - minVal)) * float64(drawArea.Dy()-1))
				point := image.Pt(drawArea.Min.X+(int(math.Round(float64(j)*self.HorizontalScale))), drawArea.Max.Y-1-height)
				if point.In(drawArea) {
					buf.SetCell(
						NewCell(self.DotMarkerRune, NewStyle(SelectColor(self.LineColors, i))),
						point,
					)
				}
			}
		}
	case LineChart:
		for i, line := range self.Data {
			for j := 0; j < len(line) && int(math.Round(float64(j)*self.HorizontalScale)) < drawArea.Dx(); j++ {
				val := line[j]
				height := int((val / (maxVal - minVal)) * float64(drawArea.Dy()-1))
				buf.SetCell(
					NewCell(self.DotMarkerRune, NewStyle(SelectColor(self.LineColors, i))),
					image.Pt(drawArea.Min.X+(int(math.Round(float64(j)*self.HorizontalScale))), drawArea.Max.Y-1-height),
				)
			}
		}
	}
}

func (self *Plot) plotAxes(buf *Buffer, maxVal, minVal float64) {
	// draw origin cell
	buf.SetCell(
		NewCell(BOTTOM_LEFT, NewStyle(self.AxesColor)),
		image.Pt(self.Inner.Min.X+self.yAxisLabelsWidth, self.Inner.Max.Y-xAxisLabelsHeight-1),
	)
	// draw x axis line
	for i := self.yAxisLabelsWidth + 1; i < self.Inner.Dx(); i++ {
		buf.SetCell(
			NewCell(HORIZONTAL_DASH, NewStyle(self.AxesColor)),
			image.Pt(i+self.Inner.Min.X, self.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}
	// draw y axis line
	for i := 0; i < self.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			NewCell(VERTICAL_DASH, NewStyle(self.AxesColor)),
			image.Pt(self.Inner.Min.X+self.yAxisLabelsWidth, i+self.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		"0",
		self.LabelStyle,
		image.Pt(self.Inner.Min.X+self.yAxisLabelsWidth, self.Inner.Max.Y-1),
	)
	// draw rest
	if self.HorizontalScale == 0 {
		return
	}
	// TODO: use floats everywhere and just convert back to int when done
	for x := self.Inner.Min.X + self.yAxisLabelsWidth + int(math.Round(float64(xAxisLabelsGap)*self.HorizontalScale)) + 1; x < self.Inner.Max.X-1; {
		label := fmt.Sprintf(
			"%d",
			int((float64(x)-float64(self.Inner.Min.X+self.yAxisLabelsWidth)-1)/(self.HorizontalScale)+1),
		)
		buf.SetString(
			label,
			self.LabelStyle,
			image.Pt(x, self.Inner.Max.Y-1),
		)
		x += int(float64((len(label) + xAxisLabelsGap)) * self.HorizontalScale)
	}
	// draw y axis labels
	verticalScale := (maxVal - minVal) / float64(self.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < self.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf(self.YAxisFormat, (float64(i))*verticalScale*(yAxisLabelsGap+1)+minVal),
			self.LabelStyle,
			image.Pt(self.Inner.Min.X, self.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (self *Plot) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = GetMaxFloat64From2dSlice(self.Data)
	}
	minVal := self.MinVal
	if minVal == 0 {
		minVal, _ = GetMinFloat64From2dSlice(self.Data)
	}

	if self.ShowAxes {
		self.yAxisLabelsWidth = len(fmt.Sprintf(self.YAxisFormat, maxVal))
		self.plotAxes(buf, maxVal, minVal)
	}

	drawArea := self.Inner
	if self.ShowAxes {
		drawArea = image.Rect(
			self.Inner.Min.X+self.yAxisLabelsWidth+1, self.Inner.Min.Y,
			self.Inner.Max.X, self.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	switch self.Marker {
	case MarkerBraille:
		self.renderBraille(buf, drawArea, maxVal, minVal)
	case MarkerDot:
		self.renderDot(buf, drawArea, maxVal, minVal)
	}
}
