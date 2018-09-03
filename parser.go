package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func parseFile(srcPath string) ([]TextBlock, error) {
	src, err := ioutil.ReadFile(filepath.Join(srcPath))
	if err != nil {
		return nil, err
	}

	tbs := make([]TextBlock, 0)

	brRgx := regexp.MustCompile(`\((.*?)\)`)
	singleVerseRgx := regexp.MustCompile(`(\p{L}{2,3})\s?.?(\d+):(\d+)`)
	singleChapterRgx := regexp.MustCompile(`(\p{L}{2,3})\s?.?(\d+):(\d+)-(\d+)`)
	multiChapterRgx := regexp.MustCompile(`(\p{L}{2,3})\s?\.?(\d+):(\d+)-(\d+):(\d+)`)
	for _, rawSub := range brRgx.FindAllStringSubmatch(string(src), -1) {
		raw := rawSub[1]
		// Find all multi chapter passages
		for _, r := range multiChapterRgx.FindAllStringSubmatch(raw, -1) {
			cns, _ := strconv.Atoi(r[2])
			vns, _ := strconv.Atoi(r[3])
			cne, _ := strconv.Atoi(r[4])
			vne, _ := strconv.Atoi(r[5])
			tb := TextBlock{
				BookName:     strings.ToLower(strings.TrimSpace(strings.Replace(r[1], ".", "", -1))),
				StartChapter: cns - 1,
				StartVerse:   vns - 1,
				EndChapter:   cne - 1,
				EndVerse:     vne - 1,
			}
			tbs = append(tbs, tb)
		}

		// Remove already processed
		raw = multiChapterRgx.ReplaceAllString(raw, "")

		// Find all single chapter multi verses passages
		for _, r := range singleChapterRgx.FindAllStringSubmatch(raw, -1) {
			cns, _ := strconv.Atoi(r[2])
			vns, _ := strconv.Atoi(r[3])
			vne, _ := strconv.Atoi(r[4])
			tb := TextBlock{
				BookName:     strings.ToLower(strings.TrimSpace(strings.Replace(r[1], ".", "", -1))),
				StartChapter: cns - 1,
				StartVerse:   vns - 1,
				EndChapter:   cns - 1,
				EndVerse:     vne - 1,
			}
			tbs = append(tbs, tb)
		}

		// Remove already processed
		raw = singleChapterRgx.ReplaceAllString(raw, "")

		// Find all single chapter multi verses passages
		for _, r := range singleVerseRgx.FindAllStringSubmatch(raw, -1) {
			cn, _ := strconv.Atoi(r[2])
			vn, _ := strconv.Atoi(r[3])
			tb := TextBlock{
				BookName:     strings.ToLower(strings.TrimSpace(strings.Replace(r[1], ".", "", -1))),
				StartChapter: cn - 1,
				StartVerse:   vn - 1,
				EndChapter:   cn - 1,
				EndVerse:     vn - 1,
			}
			tbs = append(tbs, tb)
		}

	}

	return tbs, nil
}

func (m *Module) GetScripture(tb TextBlock) (string, error) {
	book, ok := m.Books[tb.BookName]
	if !ok {
		return "", fmt.Errorf("book %s not found in module", tb.BookName)
	}

	if tb.StartChapter > tb.EndChapter {
		return "", fmt.Errorf("worng chapters order")
	} else if tb.StartChapter == tb.EndChapter && tb.StartVerse > tb.EndVerse {
		return "", fmt.Errorf("worng verses order")
	}

	verses := make([]string, 0)
	var fv, lv int
	for c := tb.StartChapter; c <= tb.EndChapter; c++ {
		fv = 0
		if c == tb.StartChapter {
			fv = tb.StartVerse
		}
		lv = len(book.Chapters[c].Verses) - 1
		if c == tb.EndChapter {
			lv = tb.EndVerse
		}
		for v := fv; v <= lv; v++ {
			if tb.StartChapter != tb.EndChapter {
				verses = append(verses, fmt.Sprintf("%d:%s", c+1, book.Chapters[c].Verses[v]))
			} else {
			verses = append(verses, book.Chapters[c].Verses[v])
			}
		}
	}
	return strings.Join(verses, "\n"), nil
}