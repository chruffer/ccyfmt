package ccyfmt

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

const NUMBERPLACEHOLDER = "#"

type ccies struct {
	Currencies map[string]ccy
}

type ccy struct {
	Name        string
	Symbol      string
	ISO4217Code string
	ISO4217Num  string
	Majorname   string
	Minorname   string
	Decimals    string
	Format      string
}

func NewCcies() (ccies, error) {

	c := ccies{}
	err := json.Unmarshal([]byte(CCYFORMATS), &c.Currencies)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c *ccies) FormatCurrency(num string, dec string, cur string) (string, error) {

	// get the currency
	ccy := c.Currencies[cur]

	// check lenght of dec
	decInt, err := strconv.Atoi(ccy.Decimals)
	if err != nil {
		return "", err
	}
	if len(dec) > decInt {
		return "", errors.New("decimals to long for this ccy")
	}
	if len(dec) < decInt {
		return "", errors.New("not enough decimals for this ccy")
	}

	// concat num and dec
	wholeNum := num + dec
	l := len(wholeNum)

	// create specific format string for number
	fdec := identifyFormatDecimalsPart(ccy.Format, ccy.Decimals)
	funreg := identifyUnregularFormatPart(ccy.Format, ccy.Decimals)
	frep := identifyRepeatFormatPart(ccy.Format)
	formatString := funreg + fdec
	lf := len(formatString)
	missing := l - lf
	repeat := missing/(len(frep)-1) + 1 // - separator + one more
	frep = strings.Repeat(frep, repeat)
	formatString = frep + formatString
	lf = len(formatString)

	// return value
	var f string

	// make array and iterate from back to front
	n := []byte(wholeNum)
	t := []byte(formatString)

	for inum, ifor := l-1, lf-1; inum >= 0; inum-- {
		//concat the sign if it's not a numplaceholder
		if string(t[ifor]) != NUMBERPLACEHOLDER {
			f = string(t[ifor]) + f
			ifor-- // reduce index
		}
		// add the number
		f = string(n[inum]) + f
		ifor-- // reduce index

	}
	return f, nil

}

func identifyRepeatFormatPart(f string) string {
	var repeatPart string
	var start bool

	i := 0

	for _, v := range []byte(f) {

		// end if second separarator is found
		if start && string(v) != NUMBERPLACEHOLDER {
			break
		}
		// start when first separator is found
		if string(v) != NUMBERPLACEHOLDER {
			start = true
		}
		// add separator or repeated numberplaceholders
		if start {
			repeatPart += string(v)
		}
		i++
	}
	return repeatPart
}

// static part is needed for currencies like INR which have an
// unregular formatpart at the first digets 1,23,000.00
// this will return ,000
func identifyUnregularFormatPart(f string, decimals string) string {

	// cutoff the dec part
	f = strings.TrimSuffix(f, identifyFormatDecimalsPart(f, decimals))

	var unregularPart string
	var separatorCount = 0

	// second separator indicates unregular formatpart

	for _, v := range []byte(f) {
		// set active when first separator is found
		if string(v) != NUMBERPLACEHOLDER {
			separatorCount++
		}
		// add chars when second separator is active
		if separatorCount == 2 {
			unregularPart += string(v)
		}
	}
	return unregularPart
}

// mirror the format string
func reverseFormat(f string) string {
	// reverse format
	l := len(f)
	rf := make([]byte, l)
	i := 1
	for _, v := range f {

		rf[l-i] = byte(v)
		i++
	}
	return string(rf)
}

// get the decimal part of the format only
func identifyFormatDecimalsPart(f string, decimals string) string {

	decimalsPart := ""

	if decimals == "0" {
		return decimalsPart
	}

	rf := reverseFormat(f)

	for _, v := range []byte(rf) {
		decimalsPart += string(v)
		if string(v) != NUMBERPLACEHOLDER {
			break
		}
	}
	decimalsPart = reverseFormat(decimalsPart)
	return decimalsPart
}
