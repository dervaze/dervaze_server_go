package lang

import (
	"errors"
	"regexp"
	"strings"

	"log"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// MakeRoot builds a Root from Latin and Visenc spelling of a word by automatically filling other information
func MakeRoot(latin string, visenc string, pos PartOfSpeech) Root {

	if visenc == "" {
		log.Println("Empty visenc for latin", latin)
	}

	ow, err := MakeOttomanWord(visenc, "")

	if err != nil {
		log.Println(err)
	}

	r := Root{
		TurkishLatin:       latin,
		Ottoman:            ow,
		LastVowel:          LastVowel(latin),
		LastConsonant:      LastConsonant(latin),
		EffectiveLastVowel: EffectiveLastVowel(latin),
		Abjad:              ow.Abjad,
		PartOfSpeech:       pos,
		EndsWithVowel:      EndsWithVowel(latin),
		HasSingleVowel:     HasSingleVowel(latin),
		LastVowelHard:      LastVowelHard(latin),
		LastConsonantHard:  LastConsonantHard(latin),
		// following three will be updated
		EffectiveTurkishLatin: latin,
		EffectiveVisenc:       ow.Visenc,
		HasConsonantSoftening: false,
	}

	UpdateEffectiveSoftening(&r)

	return r
}

// MakeOttomanWord builds an OttomanWord from either visenc or unicode
func MakeOttomanWord(visenc string, unicode string) (*OttomanWord, error) {
	if visenc == "" && unicode == "" {
		return nil, errors.New("Need either visenc or ottoman")
	}

	var cleanVisenc string
	if len(visenc) == 0 {
		cleanVisenc = UnicodeToVisenc(unicode)
	} else {
		cleanVisenc = regexp.MustCompile("[^a-z0-9 |<>]+").ReplaceAllLiteralString(visenc, "")
		if cleanVisenc != visenc {
			log.Printf("Cleaned Visenc %s -> %s", visenc, cleanVisenc)
		}
	}

	var normalized string

	if len(unicode) == 0 {
		normalized = norm.NFKC.String(VisencToUnicode(cleanVisenc))
	} else {
		normalized = norm.NFKC.String(unicode)
	}

	if !utf8.ValidString(normalized) {
		log.Printf("Invalid UTF-8 for Unicode: %s", normalized)
	}

	if !utf8.ValidString(cleanVisenc) {
		log.Printf("Invalid UTF-8 for Visenc: %s", cleanVisenc)
	}

	abjad := VisencToAbjad(cleanVisenc)

	searchKey := SearchKey(cleanVisenc)

	dotlessSearchKey := DotlessSearchKey(cleanVisenc)

	return &OttomanWord{
		Visenc:           cleanVisenc,
		Unicode:          normalized,
		Abjad:            abjad,
		SearchKey:        searchKey,
		DotlessSearchKey: dotlessSearchKey,
	}, nil
}

// SearchKey removes non letter diacritics from visenc string
func SearchKey(s string) string {
	sk := regexp.MustCompile(`([oui][0456789]+)`)
	return sk.ReplaceAllLiteralString(s, "")
}

// DotlessSearchKey removes all dots and signs from visenc string
func DotlessSearchKey(s string) string {
	sk := regexp.MustCompile(`([oui][0123456789]+)`)
	return sk.ReplaceAllLiteralString(s, "")
}

// VisencToUnicode converts a visenc string to unicode representation
func VisencToUnicode(s string) string {
	visenc := SplitVisenc(s, false)

	out := ""
	for _, v := range visenc {
		out += visencToUnicode[v]
	}
	return out
}

// UnicodeToVisenc converts a unicode string to visenc representation
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
	cleaned := SearchKey(s)
	visencLetters := SplitVisenc(cleaned, false)
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

// SplitVisenc splits s and returns letter groups according to visencToUnicode keys
func SplitVisenc(s string, addInvalidLetters bool) []string {
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

// ContainsArabicChars checks whether the string contains unicode runes between 0600 and 06FF
func ContainsArabicChars(s string) bool {
	for _, r := range s {
		if r >= 0x0600 && r <= 0x06FF {
			return true
		}
	}
	return false
}

// ContainsDigits returns true if `s` contains any digits
func ContainsDigits(s string) bool {
	for _, r := range s {
		if r >= 0x0030 && r <= 0x0039 {
			return true
		}
	}
	return false
}

// EndsWithVowel Checks whether a string ends with a vowel
func EndsWithVowel(s string) bool {
	return endsWithVowelRegex.MatchString(s)
}

// HasSingleVowel checks whether a string has a single vowel
func HasSingleVowel(s string) bool {
	return hasSingleVowelRegex.MatchString(s)
}

// LastConsonantHard checks whether a word has a final "fstkçşhp"
func LastConsonantHard(s string) bool {
	return lastConsonantHardRegex.MatchString(s)
}

// LastVowelHard checks whether a word ends with one of aıou
func LastVowelHard(s string) bool {
	ev := EffectiveLastVowel(s)
	if ev == "a" || ev == "ı" || ev == "o" || ev == "u" {
		return true
	}
	return false
}

// EffectiveLastVowel checks a word agains effectiveLastVowelRegexes to determine the vowel that governs vowel harmonization rules
func EffectiveLastVowel(s string) string {
	for r, v := range effectiveLastVowelRegexes {
		if r.MatchString(s) {
			return v
		}
	}
	return LastVowel(s)
}

// LastVowel returns the last vowel of a word
func LastVowel(s string) string {
	return ultimateVowelRegex.FindString(s)
}

// LastConsonant returns the last consonant of a word
func LastConsonant(s string) string {
	return ultimateConsonantRegex.FindString(s)
}

// UpdateEffectiveSoftening updates the EffectiveTurkishLatin, EffectiveVisenc and HasConsonantSoftening by checking suffixes for spelling
func UpdateEffectiveSoftening(r *Root) {

	if strings.HasSuffix(r.TurkishLatin, "k") &&
		strings.HasSuffix(r.Ottoman.Visenc, "fo2") {
		tlr := []rune(r.TurkishLatin)
		tll := len(tlr)

		r.EffectiveTurkishLatin = string(tlr[0:tll-1]) + "ğ"
		ovl := len(r.Ottoman.Visenc)
		r.EffectiveVisenc = r.Ottoman.Visenc[0:ovl-3] + "ao1"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "p") && strings.HasSuffix(r.Ottoman.Visenc, "bu1") {
		tlr := []rune(r.TurkishLatin)
		tll := len(tlr)
		r.EffectiveTurkishLatin = string(tlr[0:tll-1]) + "b"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "ç") && strings.HasSuffix(r.Ottoman.Visenc, "xu1") {
		tlr := []rune(r.TurkishLatin)
		tll := len(tlr)
		r.EffectiveTurkishLatin = string(tlr[0:tll-1]) + "c"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "t") && strings.HasSuffix(r.Ottoman.Visenc, "d") {
		tlr := []rune(r.TurkishLatin)
		tll := len(tlr)
		r.EffectiveTurkishLatin = string(tlr[0:tll-1]) + "d"
		r.HasConsonantSoftening = true
	}

}
