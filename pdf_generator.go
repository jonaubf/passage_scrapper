package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
	"path/filepath"
	"strings"
)

const (
	Width          = 512
	Height         = 384
	HeaderFontSize = 52
	TextFontSize   = 64
	TopMargin      = 15
	SideMargin     = 20
)

func generatePDF(m *Module, tbs []TextBlock, bgrndPath string) {
	pdf := initPDF()
	for _, tb := range tbs {
		verses, err := m.GetScripture(tb)
		if err != nil {
			fmt.Printf("[ERROR] can't get verses for %s\n", tb.String())
			continue
		}
		for {
			verses = addPage(pdf, tb.String(), bgrndPath, verses)
			if len(verses) == 0 {
				break
			}
		}
	}
	err := pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		fmt.Printf("[WARNING] error on pdf creation: %s\n", err)
	}
}

func initPDF() *gofpdf.Fpdf {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size:    gofpdf.SizeType{Wd: Width, Ht: Height},
	})
	pdf.SetFontLocation("static/")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.SetAutoPageBreak(false, 0)
	return pdf
}

func addPage(pdf *gofpdf.Fpdf, header, background string, verses []string) []string {
	pdf.AddPage()
	pdf.SetMargins(SideMargin, TopMargin, SideMargin)
	pdf.SetFont("Helvetica", "", HeaderFontSize)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	_, lineHt := pdf.GetFontSize()
	if len(background) != 0 {
		fi, err := os.Stat(background)
		assertError(err)

		imageType := "PNG"
		switch mode := fi.Mode(); {
		case mode.IsRegular():
			if ext := strings.ToLower(filepath.Ext(background)); ext == ".jpg" || ext == ".jpeg" {
				imageType = "JPG"
			} else if ext == ".png" {
				imageType = "PNG"
			} else {
				assertError(fmt.Errorf("unsupported image type"))
			}

		}

		l1 := pdf.AddLayer("Background", true)
		l2 := pdf.AddLayer("Text", true)
		pdf.BeginLayer(l1)
		opt := gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}
		pdf.ImageOptions(background, 0, 0, 0, 0, false, opt, 0, "")
		pdf.EndLayer()
		pdf.BeginLayer(l2)
		pdf.SetMargins(SideMargin, TopMargin, SideMargin)
		pdf.SetDrawColor(255, 255, 255)
		pdf.SetTextColor(255, 255, 255)
	} else {
		pdf.SetDrawColor(0, 0, 0)
		pdf.SetTextColor(0, 0, 0)
	}

	pdf.SetY(TopMargin)
	pdf.SetX(SideMargin)
	pdf.Cell(Width-(2*SideMargin), lineHt, tr(header))
	pdf.Ln(-1)
	vShift := float64(TopMargin) / 2
	pdf.SetLineWidth(2)
	pdf.Line(0, lineHt+vShift+float64(TopMargin), Width, lineHt+vShift+float64(TopMargin))

	pdf.SetY(lineHt + (2 * TopMargin))
	pdf.SetFont("Helvetica", "", TextFontSize)
	_, lineHt = pdf.GetFontSize()

	for _, v := range verses {
		lines := pdf.SplitLines([]byte(tr(v)), Width-(2*SideMargin))
		cellHt := lineHt * float64(len(lines))
		if pdf.GetY()+cellHt > Height-TopMargin {
			return verses
		}
		for _, l := range lines {
			pdf.Cell(Width-(2*SideMargin), lineHt, fmt.Sprintf("%s", l))
			pdf.Ln(-1)
		}
		pdf.SetY(pdf.GetY() + vShift)
		if len(verses) > 1 {
			verses = verses[1:]
		}
	}

	if len(background) != 0 {
		pdf.EndLayer()
	}

	return nil
}
