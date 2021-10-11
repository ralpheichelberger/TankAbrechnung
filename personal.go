package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const LOHN_ART = "0580"

func (p *Personal) loadPersonal(line string) error {
	if len(line) < 9 {
		return fmt.Errorf("personal CSV has wrong fromat")
	}
	l := strings.Split(line, ";")
	if len(l) != 3 {
		return fmt.Errorf("personal CSV has wrong fromat")
	}
	p.Personalnummer = l[0]
	p.Vorname = l[1]
	p.Nachname = l[2]
	return nil
}

func (m *MEMDB) LoadPersonalDB(inputFielPersonal *string) error {
	data, err := ioutil.ReadFile(*inputFilePersonal)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\r\n")
	for n, line := range lines {
		if n > 0 {
			p := Personal{}
			err := p.loadPersonal(line)
			if err != nil {
				break
			}
			m.InsertPersonal(p.Personalnummer, p.Vorname, p.Nachname)
		}
	}
	return nil
}
