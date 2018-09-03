package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

func generatePDF(m *Module, tbs []TextBlock) {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr:    "mm",
		Size:       gofpdf.SizeType{Wd: 400, Ht: 300},
	})
	pdf.SetFontLocation("static/")
	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	template := pdf.CreateTemplate(func(tpl *gofpdf.Tpl) {
		tpl.Image("static/2.jpg", 0, 0, 400, 300, false, "", 0, "")
	})
	_, tplSize := template.Size()
	pdf.UseTemplateScaled(template, gofpdf.PointType{X: 0, Y: 0}, tplSize.ScaleToHeight(300))

	for _, tb := range tbs {
		text, err := m.GetScripture(tb)
		if err != nil {
			fmt.Printf("[WARNING] %s\n", err)
			continue
		}
		pdf.AddPage()
		pdf.SetFont("Helvetica", "", 48)
		pdf.SetMargins(10, 40, 10)
		pdf.SetAutoPageBreak(true, 10)

		tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
		//pdf.SetTitle(tb.String(), true)

		pdf.Write(36, tr(text))

		//fmt.Println(text)
	}

	err := pdf.OutputFileAndClose("result.pdf")
	if err != nil {
		fmt.Printf("[WARNING] error on pdf creation: %s\n", err)
	}
}