// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

//go:build !windows
// +build !windows

package termui

const (
	TOP_LEFT           = '┌'
	TOP_RIGHT          = '┐'
	BOTTOM_LEFT        = '└'
	BOTTOM_RIGHT       = '┘'
	TOP_LEFT_ROUND     = '╭'
	TOP_RIGHT_ROUND    = '╮'
	BOTTOM_LEFT_ROUND  = '╰'
	BOTTOM_RIGHT_ROUND = '╯'

	VERTICAL_LINE   = '│'
	HORIZONTAL_LINE = '─'

	VERTICAL_LEFT   = '┤'
	VERTICAL_RIGHT  = '├'
	HORIZONTAL_UP   = '┴'
	HORIZONTAL_DOWN = '┬'

	QUOTA_LEFT  = '«'
	QUOTA_RIGHT = '»'

	VERTICAL_DASH   = '┊'
	HORIZONTAL_DASH = '┈'
)
