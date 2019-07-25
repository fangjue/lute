// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"strconv"
)

type ListType int

const (
	ListTypeBullet  = 0
	ListTypeOrdered = 1
)

type ListData struct {
	typ          ListType
	tight        bool
	bulletChar   string
	start        int
	delimiter    string
	padding      int
	markerOffset int
}

type List struct {
	*BaseNode
	*ListData
}

func (list *List) CanContain(nodeType NodeType) bool {
	return NodeListItem == nodeType
}

func (list *List) Finalize() {
	var item = list.firstChild

	for nil != item {
		// check for non-final list item ending with blank line:
		if endsWithBlankLine(item) && nil != item.Next() {
			list.tight = false
			break
		}

		// recurse into children of list item, to see if there are
		// spaces between any of them:
		var subitem = item.FirstChild()
		for nil != subitem {
			if endsWithBlankLine(subitem) &&
				(nil != item.Next() || nil != subitem.Next()) {
				list.tight = false
				break
			}
			subitem = subitem.Next()
		}
		item = item.Next()
	}
}

// Parse a list marker and return data on the marker (type,
// start, delimiter, bullet character, padding) or null.
func (t *Tree) parseListMarker(container Node) *ListData {
	if t.context.indent >= 4 {
		return nil
	}
	tokens := t.context.currentLine[t.context.nextNonspace:]
	data := &ListData{
		typ:          ListTypeBullet,
		tight:        true, // lists are tight by default
		markerOffset: t.context.indent,
	}

	marker := tokens[0]
	markerLength := 1
	if itemPlus == marker.typ || itemHyphen == marker.typ || itemAsterisk == marker.typ {
		data.typ = ListTypeBullet
		data.bulletChar = marker.val
	} else if marker.isNumInt() && 9 >= len(marker.val) && container.Type() != NodeParagraph {
		data.typ = ListTypeOrdered
		data.start, _ = strconv.Atoi(marker.val)
		markerLength = 2
		if itemDot == tokens[1].typ {
			data.delimiter = "."
		} else if itemCloseParen == tokens[1].typ {
			data.delimiter = ")"
		} else {
			return nil
		}
	} else {
		return nil
	}

	// make sure we have spaces after
	nextc := t.context.currentLine[t.context.nextNonspace+markerLength]
	if !nextc.isSpaceOrTab() {
		return nil
	}

	// if it interrupts paragraph, make sure first line isn't blank
	if container.Type() == NodeParagraph && t.context.currentLine[t.context.nextNonspace+markerLength:].isBlankLine() {
		return nil
	}

	// we've got a match! advance offset and calculate padding
	t.context.advanceNextNonspace()             // to start of marker
	t.context.advanceOffset(markerLength, true) // to end of marker
	spacesStartCol := t.context.column
	spacesStartOffset := t.context.offset
	for {
		t.context.advanceOffset(1, true)
		nextc = t.context.currentLine.peek(t.context.offset)
		if t.context.column-spacesStartCol >= 5 || !nextc.isSpaceOrTab() {
			break
		}
	}

	var blank_item = nil == t.context.currentLine.peek(t.context.offset)
	var spaces_after_marker = t.context.column - spacesStartCol
	if spaces_after_marker >= 5 || spaces_after_marker < 1 || blank_item {
		data.padding = markerLength + 1
		t.context.column = spacesStartCol
		t.context.offset = spacesStartOffset
		if t.context.currentLine.peek(t.context.offset).isSpaceOrTab() {
			t.context.advanceOffset(1, true)
		}
	} else {
		data.padding = markerLength + spaces_after_marker
	}
	return data
}
