package main

import (
	"flag"
	"log"
	"strings"
)

var inputFileTankAbr *string
var inputFilePersonal *string
var outputDir *string

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

	err := MemDB.LoadPersonalDB(inputFilePersonal)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = MemDB.LoadTankabrDB(inputFileTankAbr)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = MemDB.CreateAllPDF(outputDir)
	if err != nil {
		log.Fatal(err.Error())
	}
}
