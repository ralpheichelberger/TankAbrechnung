package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func (l *ALine) Load(line string) error {
	var err error
	if len(line) < 500 {
		return fmt.Errorf("length of line(=%d) less thant 500", len(line))
	}
	l.F1 = strings.TrimSpace(line[:46])
	line = line[46:]
	l.Kartennummer = strings.TrimSpace(line[:21])
	line = line[21:]
	l.Belegnummer = strings.TrimSpace(line[:10])
	l.Zeitstempel = line[10:22]
	line = line[22:]
	l.Personalnummer = strings.TrimSpace(line[:31])
	line = line[31:]
	l.F3 = strings.TrimSpace(line[:25])
	line = line[25:]
	l.F4 = strings.TrimSpace(line[:18])
	line = line[18:]
	l.F5 = strings.TrimSpace(line[:20])
	line = line[20:]
	l.F6 = strings.TrimSpace(line[:32])
	line = line[32:]
	l.F7 = strings.TrimSpace(line[:32])
	line = line[32:]
	l.F8 = strings.TrimSpace(line[:12])
	line = line[12:]
	l.Ort = strings.TrimSpace(line[:32])
	line = line[32:]
	l.F10 = strings.TrimSpace(line[:42])
	line = line[42:]
	l.Waerung1 = strings.TrimSpace(line[:3])
	l.Waerung2 = strings.TrimSpace(line[3:6])
	l.Menge, err = strconv.Atoi(strings.TrimSpace(line[6:16]))
	line = line[17:]
	l.EPreis, err = strconv.Atoi(strings.TrimSpace(line[:14]))
	l.EURBrutto, err = strconv.Atoi(strings.TrimSpace(line[14:26]))
	l.F11 = line[26:32]
	l.Rabatt, err = strconv.Atoi(strings.TrimSpace(line[32:44]))
	l.Steuersatz, err = strconv.Atoi(strings.TrimSpace(line[44:49]))
	l.Umsatzsteuer, err = strconv.Atoi(strings.TrimSpace(line[49:61]))
	l.Netto, err = strconv.Atoi(strings.TrimSpace(line[61:73]))
	l.BruttoBetrag, err = strconv.Atoi(strings.TrimSpace(line[73:85]))
	l.Gesamt, err = strconv.Atoi(strings.TrimSpace(line[85:97]))
	return err
}

func extractFields(line string) (ALine, error) {
	var aLine ALine = ALine{}
	err := aLine.Load(line)
	return aLine, err
}

func (s *SQLDB) LoadTankabrDB(inputFileTankAbr *string) error {
	data, err := ioutil.ReadFile(*inputFileTankAbr)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\r\n")
	// var dataLines []ALine = []ALine{}
	for n, line := range lines {
		if n > 0 {
			aLine, err := extractFields(line)
			if err != nil {
				//fmt.Printf("error parsing '%s'\n", line)
			} else {
				// dataLines = append(dataLines, aLine)
				s.insertTankAbrLine(
					aLine.Kartennummer,
					aLine.Personalnummer,
					aLine.Belegnummer,
					aLine.Zeitstempel,
					aLine.Waerung1,
					aLine.Waerung2,
					aLine.Ort,
					aLine.Menge,
					aLine.EPreis,
					aLine.EURBrutto,
					aLine.Steuersatz,
					aLine.Umsatzsteuer,
					aLine.BruttoBetrag,
					aLine.Rabatt,
					aLine.Netto,
					aLine.Gesamt)
			}
		}
	}
	return nil
}
