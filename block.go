// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"sync"
)

// Block is the base struct inherited by most widgets.
// Block manages size, position, border, and title.
// It implements all 3 of the methods needed for the `Drawable` interface.
// Custom widgets will override the Draw method.
type Block struct {
	Border      bool
	BorderRound bool
	BorderStyle Style

	BorderLeft, BorderRight, BorderTop, BorderBottom bool

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int

	image.Rectangle
	Inner image.Rectangle

	Title          string
	TitleStyle     Style
	TitleAlignment Alignment
	ShowTitle      bool

	// Dirty is a bool to track whether or not unrendered changes have been made
	// to a block--it is up to the user to manage this
	Dirty bool

	sync.Mutex
}

func NewBlock() *Block {
	return &Block{
		Border:       true,
		BorderStyle:  Theme.Block.Border,
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		TitleStyle: Theme.Block.Title,
		ShowTitle:  true,
	}
}

func (self *Block) drawBorder(buf *Buffer) {
	verticalCell := Cell{VERTICAL_LINE, self.BorderStyle}
	horizontalCell := Cell{HORIZONTAL_LINE, self.BorderStyle}

	// draw lines
	if self.BorderTop {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Min.Y, self.Max.X, self.Min.Y+1))
	}
	if self.BorderBottom {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Max.Y-1, self.Max.X, self.Max.Y))
	}
	if self.BorderLeft {
		buf.Fill(verticalCell, image.Rect(self.Min.X, self.Min.Y, self.Min.X+1, self.Max.Y))
	}
	if self.BorderRight {
		buf.Fill(verticalCell, image.Rect(self.Max.X-1, self.Min.Y, self.Max.X, self.Max.Y))
	}

	// draw corners
	if self.BorderTop && self.BorderLeft {
		if self.BorderRound {
			buf.SetCell(Cell{TOP_LEFT_ROUND, self.BorderStyle}, self.Min)
		} else {
			buf.SetCell(Cell{TOP_LEFT, self.BorderStyle}, self.Min)
		}
	}
	if self.BorderTop && self.BorderRight {
		if self.BorderRound {
			buf.SetCell(Cell{TOP_RIGHT_ROUND, self.BorderStyle}, self.Min)
		} else {
			buf.SetCell(Cell{TOP_RIGHT, self.BorderStyle}, image.Pt(self.Max.X-1, self.Min.Y))
		}
	}
	if self.BorderBottom && self.BorderLeft {
		if self.BorderRound {
			buf.SetCell(Cell{BOTTOM_LEFT_ROUND, self.BorderStyle}, self.Min)
		} else {
			buf.SetCell(Cell{BOTTOM_LEFT, self.BorderStyle}, image.Pt(self.Min.X, self.Max.Y-1))
		}
	}
	if self.BorderBottom && self.BorderRight {
		if self.BorderRound {
			buf.SetCell(Cell{BOTTOM_RIGHT_ROUND, self.BorderStyle}, self.Min)
		} else {
			buf.SetCell(Cell{BOTTOM_RIGHT, self.BorderStyle}, self.Max.Sub(image.Pt(1, 1)))
		}
	}
}

func (self *Block) GetTitle() string {
	return self.Title
}

func (self *Block) IsDirty() bool {
	return self.Dirty
}

func (self *Block) Clean() {
	self.Dirty = false
}

// Draw implements the Drawable interface.
func (self *Block) Draw(buf *Buffer) {
	if self.Border {
		self.drawBorder(buf)
	}
	if !self.ShowTitle {
		return
	}

	width := self.Dx()

	titleCells := ParseStyles(self.Title, self.TitleStyle)
	if len(titleCells) > width-4 {
		titleCells = titleCells[:width-3]
		titleCells[width-4] = Cell{ELLIPSES, self.TitleStyle}
	}

	switch self.TitleAlignment {
	case AlignLeft:
		buf.SetCells(
			titleCells,
			image.Pt(self.Min.X+2, self.Min.Y),
		)
	case AlignCenter:
		buf.SetCells(
			titleCells,
			image.Pt(self.Min.X+(width-len(titleCells))/2, self.Min.Y),
		)
	case AlignRight:
		buf.SetCells(
			titleCells,
			image.Pt(self.Max.X-len(titleCells)-2, self.Min.Y),
		)
	}
}

// SetRect implements the Drawable interface.
func (self *Block) SetRect(x1, y1, x2, y2 int) {
	self.Rectangle = image.Rect(x1, y1, x2, y2)
	self.Inner = image.Rect(
		self.Min.X+1+self.PaddingLeft,
		self.Min.Y+1+self.PaddingTop,
		self.Max.X-1-self.PaddingRight,
		self.Max.Y-1-self.PaddingBottom,
	)
}

// GetRect implements the Drawable interface.
func (self *Block) GetRect() image.Rectangle {
	return self.Rectangle
}
