// Package cnpj contains validation and generation of CNPJ
package cnpj

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	cnpjValidLength = 14
	invalidLength   = "CNPJ deve ter 14 dígitos."
	repeatedDigits  = "CNPJ não pode ser composto por números repetidos."
)

// Valid validates the CNPJ returning a boolean
func Valid(cnpj string) bool {
	isValid, err := AssertValid(cnpj)

	if err != nil {
		return false
	}
	return isValid
}

// AssertValid validates the CNPJ returning a boolean and the error if any
func AssertValid(cnpj string) (bool, error) {
	cnpj = sanitize(cnpj)

	if len(cnpj) != cnpjValidLength {
		return false, errors.New(invalidLength)
	}

	for i := 0; i <= 9; i++ {
		if cnpj == strings.Repeat(strconv.Itoa(i), cnpjValidLength) {
			return false, errors.New(repeatedDigits)
		}
	}

	return checkDigits(cnpj), nil
}

// Generate returns a random valid CNPJ
func Generate() string {
	rand.Seed(time.Now().UTC().UnixNano())

	cnpj := make([]int, 12)
	for i := 0; i < 12; i++ {
		cnpj[i] = rand.Intn(9)
	}
	checkDigit1 := computeCheckDigit(cnpj)
	cnpj = append(cnpj, checkDigit1)
	checkDigit2 := computeCheckDigit(cnpj)
	cnpj = append(cnpj, checkDigit2)

	var str string
	for _, value := range cnpj {
		str += strconv.Itoa(value)
	}
	return str
}

func sanitize(data string) string {
	data = strings.Replace(data, ".", "", -1)
	data = strings.Replace(data, "-", "", -1)
	data = strings.Replace(data, "/", "", -1)
	return data
}

func checkDigits(cnpj string) bool {
	data := strings.Split(cnpj, "")

	doc := make([]int, 12)
	for i := 0; i <= 11; i++ {
		digit, err := strconv.Atoi(data[i])
		if err != nil {
			return false
		}
		doc[i] = digit
	}

	checkDigit1 := computeCheckDigit(doc)
	doc = append(doc, checkDigit1)
	checkDigit2 := computeCheckDigit(doc)

	checkDigit1IsValid := strconv.Itoa(checkDigit1) == string(data[12])
	checkDigit2IsValid := strconv.Itoa(checkDigit2) == string(data[13])
	return checkDigit1IsValid && checkDigit2IsValid
}

func computeCheckDigit(doc []int) int {
	multipliers := []int{2, 3, 4, 5, 6, 7, 8, 9}

	var calc float64
	m := 0
	for i := len(doc) - 1; i >= 0; i-- {
		calc += float64(multipliers[m] * doc[i])
		m++
		if m >= len(multipliers) {
			m = 0
		}
	}
	mod := int(math.Mod(calc*10, 11))
	if mod == 10 {
		return 0
	}
	return mod
}
