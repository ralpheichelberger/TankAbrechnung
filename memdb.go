package main

import (
	"fmt"
	"strconv"
)

type AbrZeile struct {
	Belegnummer,
	Zeitstempel,
	Waerung1,
	Waerung2,
	Ort string
	Menge,
	EPreis,
	EURBrutto,
	Steuersatz,
	Umsatzsteuer,
	BruttoBetrag,
	Rabatt,
	Netto,
	Gesamt int
}

type TankAbr struct {
	Vorname,
	Nachname,
	Kartennummer string
	Zeilen []AbrZeile
}

type MEMDB struct {
	Tankabr map[int](TankAbr)
}

var MemDB MEMDB = MEMDB{
	Tankabr: make(map[int](TankAbr)),
}

func (m *MEMDB) InsertTankAbrLine(Kartennummer, Personalnummer, Belegnummer, Zeitstempel, Waerung1, Waerung2, Ort string, Menge, EPreis, EURBrutto, Steuersatz, Umsatzsteuer, BruttoBetrag, Rabatt, Netto, Gesamt int) error {
	persNr, err := strconv.Atoi(Personalnummer)
	if err != nil {
		return err
	}
	abr, found := m.Tankabr[persNr]
	if !found {
		return fmt.Errorf("ACHTUNG: Mitarbeiter mit Personalnummer %8d nicht angelegt", persNr)
	}
	abr.Kartennummer = Kartennummer
	z := AbrZeile{
		Belegnummer:  Belegnummer,
		Zeitstempel:  Zeitstempel,
		Waerung1:     Waerung1,
		Waerung2:     Waerung2,
		Ort:          Ort,
		Menge:        Menge,
		EPreis:       EPreis,
		EURBrutto:    EURBrutto,
		Steuersatz:   Steuersatz,
		Umsatzsteuer: Umsatzsteuer,
		BruttoBetrag: BruttoBetrag,
		Rabatt:       Rabatt,
		Netto:        Netto,
		Gesamt:       Gesamt,
	}
	abr.Zeilen = append(abr.Zeilen, z)
	m.Tankabr[persNr] = abr
	return nil
}

func (m *MEMDB) InsertPersonal(Personalnummer, Vorname, Nachname string) error {
	persNr, err := strconv.Atoi(Personalnummer)
	if err != nil {
		return err
	}
	abr, found := m.Tankabr[persNr]
	if !found {
		abr = TankAbr{
			Vorname:  Vorname,
			Nachname: Nachname,
			Zeilen:   []AbrZeile{},
		}
		m.Tankabr[persNr] = abr
	}
	return nil
}
