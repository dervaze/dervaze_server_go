package lang

import (
	"errors"
	"regexp"

	"golang.org/x/text/unicode/norm"
	"log"
	"unicode/utf8"
)

var visencToUnicode = map[string]string{
	"c":   "ء",
	"eo6": "آ",
	// "A": "آ",
	"e":   "ا",
	"eo5": "أ",
	"eu5": "إ",
	// "E": "أ",
	"bu1": "ب",
	// "B": "ب",
	"bu3": "پ",
	// "P": "پ",
	"bo2": "ت",
	// "T": "ت",
	"bo3": "ث",
	"xu1": "ج",
	// "C": "ج",
	"xu3": "چ",
	// "Ç": "چ",
	"x":   "ح",
	"xo1": "خ",
	// "X": "خ",
	"do1": "ذ",
	"d":   "د",
	"ro1": "ز",
	"r":   "ر",
	// "Z": "ز",
	"ro3": "ژ",
	// "J": "ژ",
	"s":   "س",
	"so3": "ش",
	// "S": "ش",
	// "Ş": "ش",
	"z":   "ص",
	"zo1": "ض",
	// "D": "ض",
	"t":   "ط",
	"to1": "ظ",
	"a":   "ع",
	"ao1": "غ",
	// "G": "غ",
	"fo1": "ف",
	// "F": "ف",
	"fo2": "ق",
	// "Q": "ق",
	"lo5": "ك",
	"ko5": "ك",
	"k":   "ک",
	"ko7": "گ",
	// "K": "گ",
	"ko3": "ڭ",
	// "lo5o3": "ڭ",
	"l":   "ل",
	"m":   "م",
	"bo1": "ن",
	// "N": "ن",
	"w":   "و",
	"wo5": "ؤ",
	"h":   "ه",
	// "h": "ە",
	// "h": "\u06D5",
	"ho2": "ة",
	"y":   "ی", // x6cc
	// "y": "ى", // x649
	"bu2": "ي",
	// "Y": "ي",
	"yo5":    "ئ",
	"bo5":    "ئ",
	"n0":     "۰",
	"n1":     "۱",
	"n2":     "۲",
	"n3":     "۳",
	"n4":     "۴",
	"n5":     "۵",
	"n6":     "۶",
	"n7":     "۷",
	"n8":     "۸",
	"n9":     "۹",
	"&zwnj;": "\u200C",
	"||":     "\u200C",
	"<>":     "\u200C",
	"&zwj;":  "\u200D",
	"><":     "\u200D",
	"&lrm;":  "\u200E",
	"&rlm;":  "\u200F",
	"&ls;":   "\u2028",
	"&ps;":   "\u2028",
	"&lre;":  "\u202A",
	"&rle;":  "\u202B",
	"&pdf;":  "\u202C",
	"&lro;":  "\u202D",
	"&rlo;":  "\u202D",
	"&bom;":  "\uFEFF",
	"o4":     "\u064E",
	"u4":     "\u0650",
	"o9":     "\u064F",
	"u44":    "\u064D",
	"o44":    "\u064B",
	"o99":    "\u064C",
	"o8":     "\u0651",
	"o0":     "\u0652",
	"o6":     "\u0653",
	" ":      " ",
	"bot":    "\u0679",
	"o5":     "\u0654",
	"u5":     "\u0655",
}

var unicodeToVisenc = map[string]string{
	"ء": "c",
	"آ": "eo6",
	// "آ": "A",
	"ا": "e",
	"أ": "eo5",
	"إ": "eu5",
	// "أ": "E",
	"ب": "bu1",
	// "ب": "B",
	"پ": "bu3",
	// "پ": "P",
	"ت": "bo2",
	// "ت": "T",
	"ث": "bo3",
	"ج": "xu1",
	// "ج": "C",
	"چ": "xu3",
	// "چ": "Ç",
	"ح": "x",
	"خ": "xo1",
	// "خ": "X",
	"د": "d",
	"ذ": "do1",
	"ر": "r",
	"ز": "ro1",
	// "ز": "Z",
	"ژ": "ro3",
	// "ژ": "J",
	"س": "s",
	"ش": "so3",
	// "ش": "S",
	// "ش": "Ş",
	"ص": "z",
	"ض": "zo1",
	// "ض": "D",
	"ط": "t",
	"ظ": "to1",
	"ع": "a",
	"غ": "ao1",
	// "غ": "G",
	"ف": "fo1",
	// "ف": "F",
	"ق": "fo2",
	// "ق": "Q",
	// "ك": "lo5",
	"ك": "k",
	"ک": "k",
	"گ": "ko7",
	// "گ": "K",
	"ڭ": "ko3",
	// "ڭ": "lo5o3",
	"ل": "l",
	"م": "m",
	"ن": "bo1",
	// "ن": "N",
	"و": "w",
	"ؤ": "wo5",
	"ه": "h",
	"ە": "h",
	"ة": "ho2",
	"ی": "y", // x6cc
	"ى": "y", // x649
	"ي": "bu2",
	// "ي": "Y",
	// "ئ": "yo5",
	"ئ": "bo5",
	"۰": "n0",
	"۱": "n1",
	"۲": "n2",
	"۳": "n3",
	"۴": "n4",
	"۵": "n5",
	"۶": "n6",
	"۷": "n7",
	"۸": "n8",
	"۹": "n9",
	// "\u200C": "&zwnj;",
	"\u200C": "||",
	// "\u200C": "<>",
	// "\u200D": "&zwj;",
	"\u200D": "><",
	"\u200E": "&lrm;",
	"\u200F": "&rlm;",
	"\u2028": "&ls;",
	// "\u2028": "&ps;",
	"\u202A": "&lre;",
	"\u202B": "&rle;",
	"\u202C": "&pdf;",
	"\u202D": "&lro;",
	// "\u202D": "&rlo;",
	"\uFEFF": "&bom;",
	"\u064E": "o4",
	"\u0650": "u4",
	"\u064F": "o9",
	"\u064D": "u44",
	"\u064B": "o44",
	"\u064C": "o99",
	"\u0651": "o8",
	"\u0652": "o0",
	"\u0653": "o6",
	" ":      " ",
	"\u0679": "bot",
	"\u0654": "o5",
	"\u0655": "u5",
}

var visencToAbjad = map[string]int32{
	"e":   1,
	"bu1": 2,
	"bu3": 2,
	"xu1": 3,
	"xu3": 3,
	"d":   4,
	"h":   5,
	"w":   6,
	"ro1": 7,
	"ro3": 7,
	"x":   8,
	"t":   9,
	"y":   10,
	"bu2": 10,
	"k":   20,
	"ko7": 20,
	"l":   30,
	"m":   40,
	"bo1": 50,
	"s":   60,
	"a":   70,
	"fo1": 80,
	"z":   90,
	"fo2": 100,
	"r":   200,
	"so3": 300,
	"bo2": 400,
	"bo3": 500,
	"xo1": 600,
	"do1": 700,
	"zo1": 800,
	"to1": 900,
	"ao1": 1000,
}

// MakeOttomanWord builds an OttomanWord from either visenc or unicode
func MakeOttomanWord(visenc string, unicode string) (*OttomanWord, error) {
	if visenc == "" && unicode == "" {
		return nil, errors.New("Need either visenc or ottoman")
	}

	var clean_visenc string
	if len(visenc) == 0 {
		clean_visenc = UnicodeToVisenc(unicode)
	} else {
		clean_visenc = regexp.MustCompile("[^a-z0-9 |<>]+").ReplaceAllLiteralString(visenc, "")
		if clean_visenc != visenc {
			log.Printf("Cleaned Visenc %s -> %s", visenc, clean_visenc)
		}
	}

	var normalized string

	if len(unicode) == 0 {
		normalized = norm.NFKC.String(VisencToUnicode(clean_visenc))
	} else {
		normalized = norm.NFKC.String(unicode)
	}

	if !utf8.ValidString(normalized) {
		log.Printf("Invalid UTF-8 for Unicode: %s", normalized)
	}

	if !utf8.ValidString(clean_visenc) {
		log.Printf("Invalid UTF-8 for Visenc: %s", clean_visenc)
	}

	abjad := VisencToAbjad(clean_visenc)

	searchKey := SearchKey(clean_visenc)

	dotlessSearchKey := DotlessSearchKey(clean_visenc)

	return &OttomanWord{
		Visenc:           clean_visenc,
		Unicode:          normalized,
		Abjad:            abjad,
		SearchKey:        searchKey,
		DotlessSearchKey: dotlessSearchKey,
	}, nil
}

func SearchKey(s string) string {
	sk := regexp.MustCompile(`([oui][0456789]+)`)
	return sk.ReplaceAllLiteralString(s, "")
}

func DotlessSearchKey(s string) string {
	sk := regexp.MustCompile(`([oui][0123456789]+)`)
	return sk.ReplaceAllLiteralString(s, "")
}

func VisencToUnicode(s string) string {
	visenc := SplitVisenc(s, false)

	out := ""
	for _, v := range visenc {
		out += visencToUnicode[v]
	}
	return out
}

func UnicodeToVisenc(s string) string {
	out := ""

	for _, u := range s {
		v, exists := unicodeToVisenc[string(u)]
		if exists {
			out += v
		}
	}

	return out
}

// VisencToAbjad calculates the abjad value for the given word in Visenc
func VisencToAbjad(s string) int32 {
	visencLetters := SplitVisenc(s, false)
	var out int32 = 0
	for _, v := range visencLetters {
		value, exists := visencToAbjad[v]
		if exists {
			out += value
		}
	}
	return out
}

// UnicodeToAbjad calculates the abjad value for a word given in Unicode by converting it to Visenc first
func UnicodeToAbjad(s string) int32 {
	return VisencToAbjad(UnicodeToVisenc(s))
}

// TFstring returns ifTrue or ifFalse according to condition
func TFstring(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

// TFString returns ifTrue or ifFalse according to condition
func TFint(condition bool, ifTrue, ifFalse int) int {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

// SplitVisenc splits s and returns letter groups according to visencToUnicode keys
func SplitVisenc(s string, addInvalidLetters bool) []string {
	// FIXME convert s to []rune and process that
	r := []rune(s)
	rlen := len(r)
	maxVisencLen := 5
	start := 0
	end := TFint(start+maxVisencLen < rlen, start+maxVisencLen, rlen)
	group := make([]string, 0, len(r))

	for start < rlen {
		_, exists := visencToUnicode[string(r[start:end])]
		if exists {
			group = append(group, string(r[start:end]))
			start = end
			end = TFint(start+maxVisencLen < rlen, start+maxVisencLen, rlen)
		} else {
			if end > start+1 {
				end--

			} else {
				if addInvalidLetters {
					group = append(group, string(r[start:end]))
				}
				start++
				end = TFint(start+maxVisencLen < rlen, start+maxVisencLen, rlen)
			}
		}

	}

	return group
}

var VISENC_TO_OTTOMAN = map[string]string{
	"c":   "ء",
	"eo6": "آ",
	// "A": "آ",
	"e":   "ا",
	"eo5": "أ",
	"eu5": "إ",
	// "E": "أ",
	"bu1": "ب",
	// "B": "ب",
	"bu3": "پ",
	// "P": "پ",
	"bo2": "ت",
	// "T": "ت",
	"bo3": "ث",
	"xu1": "ج",
	// "C": "ج",
	"xu3": "چ",
	// "Ç": "چ",
	"x":   "ح",
	"xo1": "خ",
	// "X": "خ",
	"do1": "ذ",
	"d":   "د",
	"ro1": "ز",
	"r":   "ر",
	// "Z": "ز",
	"ro3": "ژ",
	// "J": "ژ",
	"s":   "س",
	"so3": "ش",
	// "S": "ش",
	// "Ş": "ش",
	"z":   "ص",
	"zo1": "ض",
	// "D": "ض",
	"t":   "ط",
	"to1": "ظ",
	"a":   "ع",
	"ao1": "غ",
	// "G": "غ",
	"fo1": "ف",
	// "F": "ف",
	"fo2": "ق",
	// "Q": "ق",
	"lo5": "ك",
	"ko5": "ك",
	"k":   "ک",
	"ko7": "گ",
	// "K": "گ",
	"ko3": "ڭ",
	// "lo5o3": "ڭ",
	"l":   "ل",
	"m":   "م",
	"bo1": "ن",
	// "N": "ن",
	"w":   "و",
	"wo5": "ؤ",
	"h":   "ه",
	// "h": "ە",
	// "h": "\u06D5",
	"ho2": "ة",
	"y":   "ی", // x6cc
	// "y": "ى", // x649
	"bu2": "ي",
	// "Y": "ي",
	"yo5":    "ئ",
	"bo5":    "ئ",
	"n0":     "۰",
	"n1":     "۱",
	"n2":     "۲",
	"n3":     "۳",
	"n4":     "۴",
	"n5":     "۵",
	"n6":     "۶",
	"n7":     "۷",
	"n8":     "۸",
	"n9":     "۹",
	"&zwnj;": "\u200C",
	"||":     "\u200C",
	"<>":     "\u200C",
	"&zwj;":  "\u200D",
	"><":     "\u200D",
	"&lrm;":  "\u200E",
	"&rlm;":  "\u200F",
	"&ls;":   "\u2028",
	"&ps;":   "\u2028",
	"&lre;":  "\u202A",
	"&rle;":  "\u202B",
	"&pdf;":  "\u202C",
	"&lro;":  "\u202D",
	"&rlo;":  "\u202D",
	"&bom;":  "\uFEFF",
	"o4":     "\u064E",
	"u4":     "\u0650",
	"o9":     "\u064F",
	"u44":    "\u064D",
	"o44":    "\u064B",
	"o99":    "\u064C",
	"o8":     "\u0651",
	"o0":     "\u0652",
	"o6":     "\u0653",
	" ":      " ",
	"bot":    "\u0679",
	"o5":     "\u0654",
	"u5":     "\u0655",
}

// Checks whether the string contains unicode runes between 0600 and 06FF
func ContainsArabicChars(s string) bool {
	for _, r := range s {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}

func ContainsDigits(s string) bool {
	for _, r := range s {
		if r >= 0x0030 && r <= 0x0039 {
			return true
		}
	}
	return false
}
