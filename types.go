package main

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
	PathName   string
	FullName   string
	ShortName  []string
}
