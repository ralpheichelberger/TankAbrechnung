package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type InvoiceFont struct {
	ID       string
	Style    string
	FilePath string
}

var Fonts = []InvoiceFont{
	{ID: "dejavu", Style: "", FilePath: "DejaVuSans.ttf"},
	{ID: "dejavu", Style: "B", FilePath: "DejaVuSans-Bold.ttf"},
}

func CreateInvoice(b PDFData) (gofpdf.Pdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	for _, font := range Fonts {
		pdf.AddUTF8Font(font.ID, font.Style, font.FilePath)
	}

	pdf.SetLeftMargin(16)
	pdf.SetTopMargin(30)

	var maxWidth float64 = 175
	type Widths struct {
		d, o, b, m, l, s float64
	}
	var w Widths = Widths{
		d: 30,
		o: 58,
		b: 30,
		m: 12,
		l: 22,
		s: 23,
	}
	var h float64 = 6
	var hh float64 = 8
	var fsa float64 = 10
	var fst float64 = 14
	var fsl float64 = 9

	pdf.SetXY(16, 50)
	pdf.SetFont("dejavu", "B", fsa)
	pdf.CellFormat(30, h, fmt.Sprintf("%s %s", b.Vorname, b.Nachname), "", 1, "L", false, 0, "")
	pdf.SetFont("dejavu", "", fsa)
	pdf.CellFormat(10, h, fmt.Sprintf("Perso.Nr: %d", b.Personalnummer), "", 1, "L", false, 0, "")
	pdf.CellFormat(10, h, fmt.Sprintf("KartenNr: %s", b.Kartennummer), "", 1, "L", false, 0, "")

	pdf.SetXY(16, 90)
	z := b.Lines[0].Zeitstempel
	monat := fmt.Sprintf("%s.%s", z[4:6], z[:4])
	pdf.SetFont("dejavu", "B", fst)
	pdf.CellFormat(10, h, fmt.Sprintf("Tankabrechnung für %s", monat), "", 1, "L", false, 0, "")

	pdf.SetXY(16, 100)

	pdf.SetFont("dejavu", "B", fsl)
	pdf.CellFormat(10, h, "", "", 1, "L", false, 0, "")
	pdf.CellFormat(w.d, hh, "Datum", "B", 0, "L", false, 0, "")
	pdf.CellFormat(w.o, hh, "Ort", "B", 0, "L", false, 0, "")
	pdf.CellFormat(w.b, hh, "Belegnummer", "B", 0, "L", false, 0, "")
	pdf.CellFormat(w.m, hh, "Menge/l", "B", 0, "R", false, 0, "")
	pdf.CellFormat(w.l, hh, "L-Preis/EUR", "B", 0, "R", false, 0, "")
	pdf.CellFormat(w.s, hh, "Summe/EUR", "B", 1, "R", false, 0, "")
	pdf.CellFormat(10, h/2, "", "", 1, "L", false, 0, "")
	pdf.SetFont("dejavu", "", fsl)
	for _, l := range b.Lines {
		z := l.Zeitstempel
		datumzeit := fmt.Sprintf("%s.%s.%s %s:%s", z[6:8], z[4:6], z[:4], z[8:10], z[10:12])
		pdf.CellFormat(w.d, h, datumzeit, "", 0, "L", false, 0, "")
		pdf.CellFormat(w.o, h, l.Ort, "", 0, "L", false, 0, "")
		pdf.CellFormat(w.b, h, l.Belegnummer, "", 0, "L", false, 0, "")
		pdf.CellFormat(w.m, h, fmt.Sprintf("%.2f", l.Menge/100), "", 0, "R", false, 0, "")
		pdf.CellFormat(w.l, h, fmt.Sprintf("%.4f", l.EPreis/1000), "", 0, "R", false, 0, "")
		pdf.CellFormat(w.s, h, fmt.Sprintf("%.2f", l.EURBrutto/100), "", 1, "R", false, 0, "")
	}
	pdf.SetFont("dejavu", "B", fsl)
	pdf.CellFormat(maxWidth, h, fmt.Sprintf("Gesammt brutto: %.2f", b.BruttoSumme), "T", 1, "R", false, 0, "")

	return pdf, nil
}
