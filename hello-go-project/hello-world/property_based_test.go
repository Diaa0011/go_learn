package main

import (
	"strings"
)

var cases = []struct {
	Arabic int
	Roman  string
}{
	{Arabic: 1, Roman: "I"},
	{Arabic: 2, Roman: "II"},
	{Arabic: 3, Roman: "III"},
	{Arabic: 4, Roman: "IV"},
	{Arabic: 5, Roman: "V"},
	{Arabic: 6, Roman: "VI"},
	{Arabic: 7, Roman: "VII"},
	{Arabic: 8, Roman: "VIII"},
	{Arabic: 9, Roman: "IX"},
	{Arabic: 10, Roman: "X"},
	{Arabic: 14, Roman: "XIV"},
	{Arabic: 18, Roman: "XVIII"},
	{Arabic: 20, Roman: "XX"},
	{Arabic: 39, Roman: "XXXIX"},
	{Arabic: 40, Roman: "XL"},
	{Arabic: 47, Roman: "XLVII"},
	{Arabic: 49, Roman: "XLIX"},
	{Arabic: 50, Roman: "L"},
	{Arabic: 100, Roman: "C"},
	{Arabic: 90, Roman: "XC"},
	{Arabic: 400, Roman: "CD"},
	{Arabic: 500, Roman: "D"},
	{Arabic: 900, Roman: "CM"},
	{Arabic: 1000, Roman: "M"},
	{Arabic: 1984, Roman: "MCMLXXXIV"},
	{Arabic: 3999, Roman: "MMMCMXCIX"},
	{Arabic: 2014, Roman: "MMXIV"},
	{Arabic: 1006, Roman: "MVI"},
	{Arabic: 798, Roman: "DCCXCVIII"},
}

// func TestRomanNumerals(t *testing.T) {
// t.Run("1 gets converted to I", func(t *testing.T) {
// 	got := ConvertToRoman(1)
// 	want := "I"

// 	if got != want {
// 		t.Errorf("got %q, want %q", got, want)
// 	}
// })

// t.Run("2 gets converted to II", func(t *testing.T) {
// 	got := ConvertToRoman(2)
// 	want := "II"

// 	if got != want {
// 		t.Errorf("got %q, want %q", got, want)
// 	}
// })

// cases := []struct {
// 	Description string
// 	Arabic      int
// 	Want        string
// }{
// 	{"1 gets converted to I", 1, "I"},
// 	{"2 gets converted to II", 2, "II"},
// 	{"3 gets converted to III", 3, "III"},
// 	{"4 gets converted to IV (can't repeat more than 3 times)", 4, "IV"},
// 	{"5 gets converted to V", 5, "V"},
// 	{"9 gets converted to IX", 9, "IX"},
// 	{"10 gets converted to X", 10, "X"},
// 	{"14 gets converted to XIV", 14, "XIV"},
// 	{"18 gets converted to XVIII", 18, "XVIII"},
// 	{"20 gets converted to XX", 20, "XX"},
// 	{"39 gets converted to XXXIX", 39, "XXXIX"},
// 	{"40 gets converted to XL", 40, "XL"},
// 	{"47 gets converted to XLVII", 47, "XLVII"},
// 	{"49 gets converted to XLIX", 49, "XLIX"},
// 	{"50 gets converted to L", 50, "L"},
// }

// 	for _, test := range cases[:4] {
// 		t.Run(fmt.Sprintf("%q gets converted to %d", test.Roman, test.Arabic), func(t *testing.T) {
// 			// got := ConvertToRoman(test.Arabic)
// 			got := ConvertToArabic(test.Roman)
// 			if got != test.Arabic {
// 				t.Errorf("got %d, want %d", got, test.Arabic)
// 			}
// 		})
// 	}
// }

type RomanNumeral struct {
	Value  int
	Symbol string
}

var allRomanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic uint16) string {
	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for arabic >= uint16(numeral.Value) {
			result.WriteString(numeral.Symbol)
			arabic -= uint16(numeral.Value)
		}
	}

	return result.String()
}

func ConvertToArabic(roman string) int {
	// total := 0
	// for range roman {
	// 	total++
	// }
	// return total

	var arabic = 0

	for _, numeral := range allRomanNumerals {
		for strings.HasPrefix(roman, numeral.Symbol) {
			arabic += numeral.Value
			roman = strings.TrimPrefix(roman, numeral.Symbol)
		}
	}

	return arabic
}
