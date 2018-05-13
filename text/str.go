// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

// StrCount returns the number of runes in the specified string
func StrCount(s string) int {

	count := 0
	for range s {
		count++
	}
	return count
}

// StrFind returns the start and length of the rune at the
// specified position in the string
func StrFind(s string, pos int) (start, length int) {

	count := 0
	for index := range s {
		if count == pos {
			start = index
			count++
			continue
		}
		if count == pos+1 {
			length = index - start
			break
		}
		count++
	}
	if length == 0 {
		length = len(s) - start
	}
	return start, length
}

// StrRemove removes the rune from the specified string and position
func StrRemove(s string, col int) string {

	start, length := StrFind(s, col)
	return s[:start] + s[start+length:]
}

// StrInsert inserts a string at the specified character position
func StrInsert(s, data string, col int) string {

	start, _ := StrFind(s, col)
	return s[:start] + data + s[start:]
}

// StrPrefix returns the prefix of the specified string up to
// the specified character position
func StrPrefix(text string, pos int) string {

	count := 0
	for index := range text {
		if count == pos {
			return text[:index]
		}
		count++
	}
	return text
}
