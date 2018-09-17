package main

import (
	"fmt"
	"github.com/grokify/html-strip-tags-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	IniFile = "bibleqt.ini"
)

type dirIniLoader struct {
	Folder string
	Ini    string
}

type zipIniLoader struct {
	Folder string
	Ini    string
}

func (dl *dirIniLoader) Load() (*Module, error) {
	ini, err := ioutil.ReadFile(filepath.Join(dl.Folder, dl.Ini))
	if err != nil {
		return nil, err
	}

	m := &Module{
		encoding:    "utf-8",
		chapterSign: "<h4>",
		verseSign:   "<p>",
		Books:       make(map[string]*Book),
	}

	var (
		bi  *BookInfo
		bis = make([]*BookInfo, 0)
	)

	for _, s := range strings.Split(string(ini), "\n") {
		if !strings.Contains(s, "=") {
			continue
		}

		split := strings.Split(s, "=")
		if len(split) != 2 {
			continue
		}

		field := strings.TrimSpace(split[0])
		value := strings.TrimSpace(split[1])

		switch field {
		case "BibleShortName":
			m.Name = value
		case "BibleName":
			m.FullName = value
		case "DefaultEncoding":
			m.encoding = value
		case "ChapterSign":
			m.chapterSign = value
		case "VerseSign":
			m.verseSign = value
		case "PathName":
			bi = new(BookInfo)
			bi.PathName = value
		case "FullName":
			bi.FullName = value
		case "ShortName":
			bi.ShortName = strings.Split(strings.ToLower(value), " ")
		case "ChapterQty":
			bis = append(bis, bi)
		}
	}

	for _, bi = range bis {
		b, err := ioutil.ReadFile(filepath.Join(dl.Folder, bi.PathName))
		if err != nil {
			fmt.Printf("can't read book: %+v\n", *bi)
		}

		book := &Book{Name: bi.FullName, Chapters: make([]Chapter, 0)}
		var c = &Chapter{}

		for _, s := range strings.Split(string(b), "\n") {
			if strings.Contains(s, m.chapterSign) {
				// Add previous chapter to book
				if len(c.Verses) != 0 {
					book.Chapters = append(book.Chapters, *c)
				}

				// Create a new chapter
				c = new(Chapter)
				c.Verses = make([]string, 0)
			}
			if strings.Contains(s, m.verseSign) {
				s = strip.StripTags(s)
				c.Verses = append(c.Verses, s)
			}
		}

		for _, s := range bi.ShortName {
			m.Books[s] = book
		}
	}

	return m, nil
}

// Read BibeleQT module
func ReadBQTModule(modulePath string) (*Module, error) {
	fi, err := os.Stat(modulePath)
	assertError(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return loadDirModule(modulePath)
	case mode.IsRegular():
		return nil, fmt.Errorf("unsupported yet")
	}

	return nil, fmt.Errorf("path BQT module should be a folder or zip archieve")
}

func loadDirModule(mpath string) (*Module, error) {
	fInfos, err := ioutil.ReadDir(mpath)
	if err != nil {
		return nil, fmt.Errorf("can't read folder: %s", err)
	}

	for _, f := range fInfos {
		if strings.ToLower(f.Name()) == IniFile {
			dl := dirIniLoader{Folder: mpath, Ini: f.Name()}
			return dl.Load()
		}
	}

	return nil, fmt.Errorf("there is nothing to load here")
}
