package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SingleLine struct {
	Belegnummer,
	Zeitstempel,
	Waerung1,
	Waerung2,
	Ort string
	Menge,
	EPreis,
	Rabatt,
	Gesamt float64
}

type PDFData struct {
	Personalnummer int
	Kartennummer   string
	BruttoSumme    float64
	Vorname        string
	Nachname       string
	Lines          []SingleLine
}

func (m *MEMDB) CreateAllPDF(outputDir *string) error {
	var csvData string = "Personalnummer;EUR-Gesamtbetrag;Lohnart"
	os.Mkdir(*outputDir, 0770)
	for persNr, row := range MemDB.Tankabr {
		summe := 0
		var lines []SingleLine = make([]SingleLine, 0)
		for _, zeile := range row.Zeilen {
			lines = append(lines, SingleLine{
				Belegnummer: zeile.Belegnummer,
				Zeitstempel: zeile.Zeitstempel,
				Waerung1:    zeile.Waerung1,
				Waerung2:    zeile.Waerung2,
				Ort:         zeile.Ort,
				Menge:       float64(zeile.Menge),
				EPreis:      float64(zeile.EPreis),
				Rabatt:      float64(zeile.Rabatt),
				Gesamt:      float64(zeile.Gesamt),
			})
			summe += zeile.Gesamt
		}
		bruttoSumme := float64(summe) / 100
		var pdfData PDFData = PDFData{
			Kartennummer:   row.Kartennummer,
			Personalnummer: persNr,
			BruttoSumme:    bruttoSumme,
			Vorname:        row.Vorname,
			Nachname:       row.Nachname,
			Lines:          lines,
		}
		if pdfData.BruttoSumme > 0 {
			csvData = fmt.Sprintf("%s\r\n%d;%s;%s", csvData, pdfData.Personalnummer,
				strings.Replace(fmt.Sprintf("%.2f", pdfData.BruttoSumme), ".", ",", 1), LOHN_ART)
			A4PDF, err := CreateInvoice(pdfData)
			if err != nil {
				return err
			}
			var b bytes.Buffer = bytes.Buffer{}
			w := bufio.NewWriter(&b)
			err = A4PDF.Output(w)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/tankabr%06d.pdf", *outputDir, pdfData.Personalnummer), b.Bytes(), 0660)
			if err != nil {
				return err
			}
		}
	}
	err := ioutil.WriteFile(fmt.Sprintf("%s/tankabr.csv", *outputDir), []byte(csvData), 0660)
	if err != nil {
		return err
	}
	return nil
}
