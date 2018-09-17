package main

import (
	"fmt"
	"strings"
)

type Module struct {
	Name     string
	FullName string
	Books    map[string]*Book

	encoding    string
	chapterSign string
	verseSign   string
}

type Book struct {
	Name     string
	Chapters []Chapter
}

type Chapter struct {
	Verses []string
}

type BookInfo struct {
	PathName  string
	FullName  string
	ShortName []string
}

type TextBlock struct {
	BookName     string
	StartChapter int
	StartVerse   int
	EndChapter   int
	EndVerse     int
}

func (tb *TextBlock) String() string {
	if tb.StartChapter == tb.EndChapter && tb.StartVerse == tb.EndVerse {
		return fmt.Sprintf("%s.%d:%d", strings.Title(tb.BookName), tb.StartChapter+1, tb.StartVerse+1)
	} else if tb.StartChapter == tb.EndChapter {
		return fmt.Sprintf("%s.%d:%d-%d", strings.Title(tb.BookName), tb.StartChapter+1, tb.StartVerse+1, tb.EndVerse+1)
	}
	return fmt.Sprintf("%s.%d:%d-%d:%d", strings.Title(tb.BookName), tb.StartChapter+1, tb.StartVerse+1, tb.EndChapter+1, tb.EndVerse+1)
}
