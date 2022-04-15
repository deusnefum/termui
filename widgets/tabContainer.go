// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"fmt"

	. "github.com/deusnefum/termui/v3"
)

// TabContainer is a renderable widget which can be used to conditionally render certain tabs/views.
// TabContainer shows a list of Tab names.
// The currently selected tab can be found through the `ActiveTabIndex` field.
type TabContainer struct {
	Block
	ActiveTabIndex    int
	ActiveTabStyleStr string
	ActiveTabStyle    Style
	InactiveTabStyle  Style
	Tabs              []Drawable
	TabTitles         []string
}

func NewTabContainer(tabs ...Drawable) *TabContainer {
	return &TabContainer{
		Block:             *NewBlock(),
		Tabs:              tabs,
		ActiveTabStyleStr: "mod:reverse",
		ActiveTabStyle:    Theme.Tab.Active,
		InactiveTabStyle:  Theme.Tab.Inactive,
	}
}

func (self *TabContainer) ActiveTab() Drawable {
	return self.Tabs[self.ActiveTabIndex]
}

func (self *TabContainer) FocusLeft() {
	if self.ActiveTabIndex > 0 {
		self.ActiveTabIndex--
	} else {
		self.ActiveTabIndex = len(self.Tabs) - 1
	}
}

func (self *TabContainer) FocusRight() {
	if self.ActiveTabIndex < len(self.Tabs)-1 {
		self.ActiveTabIndex++
	} else {
		self.ActiveTabIndex = 0
	}
}

func (self *TabContainer) Draw(buf *Buffer) {
	self.Title = ""
	for i := range self.Tabs {
		left, right := "  ", "  "
		if i == self.ActiveTabIndex {
			left, right = " [[", fmt.Sprintf("]](%s) ", self.ActiveTabStyleStr)
		}

		if i < len(self.TabTitles) && self.TabTitles[i] != "" {
			self.Title += left + self.TabTitles[i] + right
		} else {
			self.Title += left + self.Tabs[i].GetTitle() + right
		}
		if i != len(self.Tabs)-1 {
			self.Title += string(VERTICAL_LINE)
		}
	}
	// do this last so we can show the updated title
	self.Block.Draw(buf)
	self.Tabs[self.ActiveTabIndex].SetRect(self.Inner.Min.X, self.Inner.Min.Y, self.Inner.Max.X, self.Inner.Max.Y)
	self.Tabs[self.ActiveTabIndex].Draw(buf)
}
