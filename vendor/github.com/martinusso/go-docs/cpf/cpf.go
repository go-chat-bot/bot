// Package cpf contains validation and generation of CPF
package cpf

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	cpfValidLength = 11
	invalidLength  = "CPF deve ter 11 dígitos."
	repeatedDigits = "CPF não pode ser composto por números repetidos."
)

// Valid validates the CPF returning a boolean
func Valid(cpf string) bool {
	isValid, err := AssertValid(cpf)
	if err != nil {
		return false
	}
	return isValid
}

// AssertValid validates the CPF returning a boolean and the error if any
func AssertValid(cpf string) (bool, error) {
	cpf = sanitize(cpf)
	if len(cpf) != cpfValidLength {
		return false, errors.New(invalidLength)
	}
	for i := 0; i <= 9; i++ {
		if cpf == strings.Repeat(strconv.Itoa(i), 11) {
			return false, errors.New(repeatedDigits)
		}
	}
	return checkDigits(cpf), nil
}

// Generate returns a random valid CPF
func Generate() string {
	rand.Seed(time.Now().UTC().UnixNano())

	data := make([]int, 9)
	for i := 0; i < 9; i++ {
		data[i] = rand.Intn(9)
	}
	checkDigit1 := computeCheckDigit(data)
	data = append(data, checkDigit1)
	checkDigit2 := computeCheckDigit(data)
	data = append(data, checkDigit2)

	var cpf string
	for _, value := range data {
		cpf += strconv.Itoa(value)
	}
	return cpf
}

func sanitize(data string) string {
	data = strings.Replace(data, ".", "", -1)
	data = strings.Replace(data, "-", "", -1)
	return data
}

func checkDigits(cpf string) bool {
	data := strings.Split(cpf, "")

	doc := make([]int, 9)
	for i := 0; i <= 8; i++ {
		digit, err := strconv.Atoi(data[i])
		if err != nil {
			return false
		}
		doc[i] = digit
	}

	checkDigit1 := computeCheckDigit(doc)
	doc = append(doc, checkDigit1)
	checkDigit2 := computeCheckDigit(doc)

	checkDigit1IsValid := strconv.Itoa(checkDigit1) == string(data[9])
	checkDigit2IsValid := strconv.Itoa(checkDigit2) == string(data[10])
	return checkDigit1IsValid && checkDigit2IsValid
}

func computeCheckDigit(data []int) int {
	var calc float64
	for i, j := 2, len(data)-1; j >= 0; i, j = i+1, j-1 {
		calc += float64(i * data[j])
	}
	mod := int(math.Mod(calc*10, 11))
	if mod == 10 {
		return 0
	}
	return mod
}
