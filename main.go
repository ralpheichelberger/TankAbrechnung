package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var inputFileTankAbr *string
var inputFilePersonal *string
var outputDir *string
var SQL_DB SQLDB

func init() {
	inputFileTankAbr = flag.String("f", "", "-f Dateiname der Abrechung")
	inputFilePersonal = flag.String("p", "", "-p Dateiname der Personalliste")
	outputDir = flag.String("o", "", "-o Ordner für PDFs")
	flag.Parse()
}

func main() {
	if *inputFileTankAbr == "" {
		log.Fatal("Bitte geben Sie den Dateinamen an der importiert werden soll\n\ttankabrechung -f <Dateiname> -o <Ordnername>\n")
	}
	if *outputDir == "" {
		log.Fatal("Bitte geben Sie ein Verzeichnis für die zuerzeugenden Abrechungen an\n\ttankabrechung -f <Dateiname> -o <Ordnername>\n")
	}
	*outputDir = strings.TrimSpace(*outputDir)
	err := SQL_DB.CreateDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer SQL_DB.db.Close()

	err = SQL_DB.CreateTables()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SQL_DB.LoadPersonalDB(inputFilePersonal)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SQL_DB.LoadTankabrDB(inputFileTankAbr)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SQL_DB.CreateAllPDF(outputDir)
	if err != nil {
		log.Fatal(err.Error())
	}
	fehlendePersonalnummern, err := SQL_DB.FindMissingPersonalnummern()
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(fehlendePersonalnummern) > 0 {
		fmt.Println("*** ACHTUNG: ")
		fmt.Println("In der Personal CSV fehlen folgende Personalnummern:")
		for _, personalNummer := range fehlendePersonalnummern {
			fmt.Println(personalNummer)
		}
		fmt.Println("Für diese Mitarbeiter konnten keine Austrucke erstellt werden!")
	}
}
