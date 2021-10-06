package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractFields(t *testing.T) {
	line := "1020210831020002986611007160001009773000000000710207026826003030   DT2609560120210803194212000314                       0000001                  NNS               AP0401107           ENI EHRWALD                     SCHANZ 1A                       6632        EHRWALD                         3202030203SUPER 95 ADD  mind.             EUREUR0000004701 0000000000131900000000620100040000000000018802000000000001003000000005010000000006201000000006013                                                     "
	var want ALine = ALine{
		Kartennummer:   "710207026826003030",
		Personalnummer: "12000314",
		Belegnummer:    "DT26095601",
		Zeitstempel:    "202108031942",
		Waerung1:       "EUR",
		Waerung2:       "EUR",
		Menge:          4701,
		EPreis:         1319,
		EURBrutto:      6201,
		Steuersatz:     2000,
		Umsatzsteuer:   1003,
		BruttoBetrag:   6201,
		Rabatt:         188,
		Netto:          5010,
		Gesamt:         6013,
		F1:             "1020210831020002986611007160001009773000000000",
		F3:             "0000001",
		F4:             "NNS",
		F5:             "AP0401107",
		F6:             "ENI EHRWALD",
		F7:             "SCHANZ 1A",
		F8:             "6632",
		Ort:            "EHRWALD",
		F10:            "3202030203SUPER 95 ADD  mind.",
		F11:            "000400",
	}
	got, err := extractFields(line)
	if assert.NoError(t, err) {
		if assert.EqualValues(t, want, got) {
			line = "1020210831053501330011007160001009773000000000710207026826002941   DT0404163920210809121712000970                       0000000                  NNS               BP9850P292                                                                                      LUBIEN KUJAWSKI                 1101030100ON ACT           LTR            PLNEUR0000004747 0000000000571000000000593800000000000000000002300000000005069000000022036000000027105000000027105                                                     "
			want = ALine{
				Kartennummer:   "710207026826002941",
				Personalnummer: "12000970",
				Belegnummer:    "DT04041639",
				Zeitstempel:    "202108091217",
				Waerung1:       "PLN",
				Waerung2:       "EUR",
				Menge:          4747,
				EPreis:         5710,
				EURBrutto:      5938,
				Steuersatz:     2300,
				Umsatzsteuer:   5069,
				BruttoBetrag:   27105,
				Rabatt:         0,
				Netto:          22036,
				Gesamt:         27105,
				F1:             "1020210831053501330011007160001009773000000000",
				F3:             "0000000",
				F4:             "NNS",
				F5:             "BP9850P292",
				F6:             "",
				F7:             "",
				F8:             "",
				Ort:            "LUBIEN KUJAWSKI",
				F10:            "1101030100ON ACT           LTR",
				F11:            "000000",
			}
			got, err = extractFields(line)
			if assert.NoError(t, err) {
				assert.EqualValues(t, want, got)
			}
		}
	}
}
