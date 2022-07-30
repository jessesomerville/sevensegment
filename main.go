package main

import (
	"fmt"
	"strings"
	"sync"
)

const (
	DefaultBlockChar = '#'
	FullBlockChar    = '█'
)

var (
	segments = []string{
		"dbcfeag",
		"cgaed",
		"fe",
		"bfgad",
		"aefcdb",
		"efa",
		"efgda",
		"gcef",
		"dcaebg",
		"dfeagc",
	}
)

func main() {
	s := NewSevenSegmentDisplay(DefaultBlockChar)
	s.PrintInt(1234567890)
}

type SevenSegment struct {
	BlockChar           string
	a, b, c, d, e, f, g bool
	lines               []string
}

func NewSevenSegmentDisplay(blockChar rune) *SevenSegment {
	return &SevenSegment{
		BlockChar: string(blockChar),
		lines:     []string{"", "", "", "", "", "", ""},
	}
}

func (s *SevenSegment) PrintInt(num int) {
	digits := newDigitsStack()
	for num > 0 {
		digits.Push(num % 10)
		num /= 10
	}
	for !digits.isEmpty {
		s.setActiveSegments(digits.Pop())
		s.horizontal(0)
		s.vertical(true)
		s.horizontal(3)
		s.vertical(false)
		s.horizontal(6)
	}

	fmt.Println(strings.Join(s.lines, "\n"))
}

func (s *SevenSegment) horizontal(lineNum int) {
	if lineNum == 3 {
		if s.d {
			s.lines[3] += " " + strings.Repeat(s.BlockChar, 4) + "  "
		} else {
			s.lines[3] += "       "
		}
		return
	}

	var mainSet, leftSet, rightSet bool
	if lineNum == 0 {
		mainSet, leftSet, rightSet = s.a, s.b, s.c
	} else {
		mainSet, leftSet, rightSet = s.g, s.e, s.f
	}

	if mainSet {
		s.lines[lineNum] += " " + strings.Repeat(s.BlockChar, 4) + "  "
	} else if leftSet && rightSet {
		s.lines[lineNum] += s.BlockChar + "    " + s.BlockChar + " "
	} else if leftSet {
		s.lines[lineNum] += s.BlockChar + "      "
	} else {
		s.lines[lineNum] += "     " + s.BlockChar + " "
	}
}

func (s *SevenSegment) vertical(topPair bool) {
	var leftSet, rightSet bool
	var lineNum int
	if topPair {
		leftSet, rightSet = s.b, s.c
		lineNum = 1
	} else {
		leftSet, rightSet = s.e, s.f
		lineNum = 4
	}

	var line string
	if leftSet && rightSet {
		line = s.BlockChar + "    " + s.BlockChar + " "
	} else if leftSet {
		line = s.BlockChar + "      "
	} else {
		line = "     " + s.BlockChar + " "
	}
	s.lines[lineNum] += line
	s.lines[lineNum+1] += line
}

func (s *SevenSegment) setActiveSegments(num int) {
	A, B, C, D := getBits(num)
	s.a = A || C || (B && D) || (!B && !D)
	s.b = A || (!C && !D) || (B && !C) || (B && !D)
	s.c = (!A && !B) || (!C && !D) || (A && D) || (C && D)
	s.d = A || (B && !(C && D)) || (!B && C)
	s.e = (!B && !D) || (C && !D) || (A && !D)
	s.f = B || D || !C
	s.g = A || (C && !D) || (B && !C && D) || (!B && C) || (!B && !D)
}

func getBits(num int) (A, B, C, D bool) {
	A = (num & 8) == 8
	B = (num & 4) == 4
	C = (num & 2) == 2
	D = (num & 1) == 1
	return A, B, C, D
}

type digitsStack struct {
	mu      sync.Mutex
	stack   []int
	isEmpty bool
}

func newDigitsStack() *digitsStack {
	return &digitsStack{isEmpty: true}
}

func (d *digitsStack) Push(digit int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.stack = append(d.stack, digit)
	d.isEmpty = false
}

func (d *digitsStack) Pop() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.stack) == 0 {
		return -1
	}
	l := len(d.stack)
	val := d.stack[l-1]
	d.stack = d.stack[:l-1]
	if l == 1 {
		d.isEmpty = true
	}
	return val
}
