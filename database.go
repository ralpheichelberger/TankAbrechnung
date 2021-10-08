package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const LOHN_ART = "0580"

type SQLDB struct {
	db *sql.DB
}

func (s *SQLDB) CreateDatabase() error {
	os.Remove("tankabrechnung.db")
	file, err := os.Create("tankabrechnung.db")
	if err != nil {
		return err
	}
	file.Close()
	sqliteDatabase, err := sql.Open("sqlite3", "./tankabrechnung.db")
	if err != nil {
		return err
	}
	s.db = sqliteDatabase
	return nil
}

func (s *SQLDB) CreateTables() error {
	createTankabr := `CREATE TABLE tankabr (
		Kartennummer string,
		Personalnummer string,
		Belegnummer    string,
		Zeitstempel    string,
		Waerung1       string,
		Waerung2       string,
		Ort            string,
		Menge          integer,
		EPreis         integer,
		EURBrutto      integer,
		Steuersatz     integer,
		Umsatzsteuer   integer,
		BruttoBetrag   integer,
		Rabatt         integer,
		Netto          integer,
		Gesamt         integer
	  );`

	statement, err := s.db.Prepare(createTankabr)
	if err != nil {
		return err
	}
	statement.Exec()
	createPersonal := `CREATE TABLE personal (
		Personalnummer string,
		Vorname string,
		Nachname string
	  );`

	log.Println("Create table 'personal'...")
	statement, err = s.db.Prepare(createPersonal)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("table created")
	return nil
}

func (s *SQLDB) insertTankAbrLine(Kartennummer, Personalnummer, Belegnummer, Zeitstempel, Waerung1, Waerung2, Ort string, Menge, EPreis, EURBrutto, Steuersatz, Umsatzsteuer, BruttoBetrag, Rabatt, Netto, Gesamt int) {
	insertLineSQL := `INSERT INTO tankabr(Kartennummer, Personalnummer, Belegnummer, Zeitstempel, Waerung1, Waerung2, Ort, Menge, EPreis, EURBrutto, Steuersatz, Umsatzsteuer, BruttoBetrag, Rabatt, Netto, Gesamt) VALUES (?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := s.db.Prepare(insertLineSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(Kartennummer, Personalnummer, Belegnummer, Zeitstempel, Waerung1, Waerung2, Ort, Menge, EPreis, EURBrutto, Steuersatz, Umsatzsteuer, BruttoBetrag, Rabatt, Netto, Gesamt)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (s *SQLDB) insertPersonal(Personalnummer, Vorname, Nachname string) {
	insertLineSQL := `INSERT INTO personal(Personalnummer, Vorname, Nachname) VALUES (?, ?, ?)`
	statement, err := s.db.Prepare(insertLineSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(Personalnummer, Vorname, Nachname)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

type PDFData struct {
	Kartennummer   string
	Personalnummer string
	BruttoSumme    float64
	Vorname        string
	Nachname       string
	Lines          []SingleLine
}

func (s *SQLDB) CreateAllPDF(outputDir *string) error {
	row, err := s.db.Query("select p.Vorname,p.Nachname,p.Personalnummer,t.Kartennummer, sum(EURBrutto) from tankabr t left join personal p on t.Personalnummer=p.Personalnummer group by t.Personalnummer;")
	if err != nil {
		log.Fatal(err)
	}
	var csvData string = "Personalnummer;EUR-Gesamtbetrag;Lohnart"
	os.Mkdir(*outputDir, 0770)
	defer row.Close()
	for row.Next() {
		var pdfData PDFData = PDFData{}
		var summe int
		row.Scan(&pdfData.Vorname, &pdfData.Nachname, &pdfData.Personalnummer, &pdfData.Kartennummer, &summe)
		pdfData.BruttoSumme = float64(summe) / 100
		csvData = fmt.Sprintf("%s\r\n%s;%s;%s", csvData, pdfData.Personalnummer,
			strings.Replace(fmt.Sprintf("%.2f", pdfData.BruttoSumme), ".", ",", 1), LOHN_ART)
		if len(pdfData.Personalnummer) == 8 {
			pdfData.Lines, err = s.GetSingleLines(pdfData.Personalnummer)
			if err != nil {
				return err
			}
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
			err = ioutil.WriteFile(fmt.Sprintf("%s/tankabr%s.pdf", *outputDir, pdfData.Personalnummer), b.Bytes(), 0660)
			if err != nil {
				return err
			}
		}
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/tankabr%s.csv", *outputDir, *outputDir), []byte(csvData), 0660)
	if err != nil {
		return err
	}
	return nil
}

type SingleLine struct {
	Belegnummer,
	Zeitstempel,
	Waerung1,
	Waerung2,
	Ort string
	Menge,
	EPreis,
	EURBrutto float64
}

func (s *SQLDB) GetSingleLines(personalNummer string) ([]SingleLine, error) {
	var lines []SingleLine = make([]SingleLine, 0)
	row, err := s.db.Query(fmt.Sprintf("select Belegnummer,Zeitstempel,Waerung1,Waerung2,Ort,Menge,EPreis,EURBrutto from tankabr t where t.Personalnummer==%s;", personalNummer))
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var l SingleLine = SingleLine{}
		row.Scan(&l.Belegnummer, &l.Zeitstempel, &l.Waerung1, &l.Waerung2, &l.Ort, &l.Menge, &l.EPreis, &l.EURBrutto)
		lines = append(lines, l)
	}
	return lines, nil
}

func (s *SQLDB) findPerson(personalNummer string) (bool, error) {
	row, err := s.db.Query(fmt.Sprintf("select Personalnummer from personal where Personalnummer='%s' limit 1", personalNummer))
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var found bool = false
	row.Next()
	var nr string
	row.Scan(&nr)
	found = len(nr) == 8
	return found, nil
}

func (s *SQLDB) FindMissingPersonalnummern() ([]string, error) {
	row, err := s.db.Query("select Personalnummer from tankabr group by Personalnummer")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var p []string = make([]string, 0)
	for row.Next() {
		var personalNummer string
		row.Scan(&personalNummer)
		found, err := s.findPerson(personalNummer)
		if err != nil {
			return p, err
		}
		if !found {
			p = append(p, personalNummer)
		}
	}
	return p, nil
}
