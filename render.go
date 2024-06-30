// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
	"sync"

	tb "github.com/nsf/termbox-go"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	GetTitle() string
	sync.Locker
}

type ConditionalDrawable interface {
	Drawable
	IsDirty() bool
	Clean()
}

func Render(items ...Drawable) {
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
	}
	tb.Flush()
}

func ConditionalRender(items ...ConditionalDrawable) {
	updateMade := false
	for _, item := range items {
		if !item.IsDirty() {
			continue
		}
		updateMade = true
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Clean()
		item.Unlock()
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
	}
	if updateMade {
		tb.Flush()
	}
}
