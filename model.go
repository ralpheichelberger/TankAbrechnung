package main

type Personal struct {
	Personalnummer,
	Vorname,
	Nachname string
}

type ALine struct {
	Kartennummer   string
	Personalnummer string
	Belegnummer    string
	Zeitstempel    string
	Waerung1       string
	Waerung2       string
	Menge          int
	EPreis         int
	EURBrutto      int
	Steuersatz     int
	Umsatzsteuer   int
	BruttoBetrag   int
	Rabatt         int
	Netto          int
	Gesamt         int
	F1             string
	F3             string
	F4             string
	F5             string
	F6             string
	F7             string
	F8             string
	Ort            string
	F10            string
	F11            string
}
