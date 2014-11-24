package ccyfmt

import (
	"testing"
)

func TestRevertFormat(t *testing.T) {

	rp := reverseFormat("#,###.##")
	if rp != "##.###,#" {
		t.Error("reversed format wrong")
	}

}

func TestNumberFormat(t *testing.T) {

	ccies, err := NewCcies()
	if err != nil {
		t.Error("error: ", err)
	}
	fs, err := ccies.FormatCurrency("19999999", "00", "EUR")
	if err != nil {
		t.Error("error: ", err)
	}
	if fs != "19.999.999,00" {
		t.Error("format wrong", fs)
	}

	fs, err = ccies.FormatCurrency("0", "99", "EUR")
	if err != nil {
		t.Error("error: ", err)
	}
	if fs != "0,99" {
		t.Error(" format wrong", fs)
	}

	fs, err = ccies.FormatCurrency("100000000000", "00", "INR")
	if err != nil {
		t.Error("error: ", err)
	}
	if fs != "1,00,00,00,00,000.00" {
		t.Error(" format wrong", fs)
	}
}

func TestNumberFormatRepeatedPart(t *testing.T) {

	rp := identifyRepeatFormatPart("#,###.##")
	if rp != ",###" {
		t.Error("repeated format wrong", rp)
	}

}

func TestNumberFormatUnregularPart(t *testing.T) {

	rp := identifyUnregularFormatPart("#,##,###.##", "2")
	if rp != ",###" {
		t.Error("unrepeated format wrong", rp)
	}

	rp = identifyUnregularFormatPart("#,###.##", "2")
	if rp != "" {
		t.Error("unrepeated format wrong", rp)
	}

}

func TestNumberDecimalsPart(t *testing.T) {

	d := identifyFormatDecimalsPart("#,###.##", "2")
	if d != ".##" {
		t.Error("decimals format wrong", d)
	}

}
