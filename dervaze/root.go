package dervaze

import (
	"deneme/dervaze/dervaze"
	"regexp"
	"strings"
)

// VOWELS are the all Roman vowels recognized
const VOWELS = "aâeıioöuûü"

var vowelRegex *regexp.Regexp = regexp.MustCompile(`[aeıioöuüâîûAEIİOÖUÜÂÎÛ]`)
var ultimateVowelRegex *regexp.Regexp = regexp.MustCompile(`.*([aâeıiîoöuüû])[^aâeıiîoöuüû]*$`)
var ultimateConsonantRegex *regexp.Regexp = regexp.MustCompile(`.*([^aâeıiîoöuüû])[aâeıiîoöuüû]*$`)
var vowelsRegex *regexp.Regexp = regexp.MustCompile(`([aâeıiîoöuüû])`)
var consonantsRegex *regexp.Regexp = regexp.MustCompile(`([bcçdfgğhjklmnpqrsştvwyxz])`)
var endsWithVowelRegex *regexp.Regexp = regexp.MustCompile(`.*[aeıioöuüâûî][']?$`)
var hasSingleVowelRegex *regexp.Regexp = regexp.MustCompile(`^[^aâeıiîoöuüû]*[aâeıiîoöuüû][^aâeıiîoöuüû]*$`)
var lastConsonantHardRegex *regexp.Regexp = regexp.MustCompile(`[fstkçşhp]'?$`)
var effectiveLastVowelRegexes = map[*regexp.Regexp]string{
	regexp.MustCompile(`.*a[^aeıioöuüâûî]*$`):     `a`,
	regexp.MustCompile(`.*â[lk][^aeıioöuüâûî]*$`): `i`,
	regexp.MustCompile(`.*â[^lkaeıioöuüâûî]*$`):   `a`,
	regexp.MustCompile(`.*e[^aeıioöuüâûî]*$`):     `e`,
	regexp.MustCompile(`.*i[^aeıioöuüâûî]*$`):     `i`,
	regexp.MustCompile(`.*î[^aeıioöuüâûî]*$`):     `i`,
	regexp.MustCompile(`.*ı[^aeıioöuüâûî]*$`):     `ı`,
	regexp.MustCompile(`.*ö[^aeıioöuüâûî]*$`):     `ö`,
	regexp.MustCompile(`.*o[^aeıioöuüâûî]*$`):     `o`,
	regexp.MustCompile(`.*ü[^aeıioöuüâûî]*$`):     `ü`,
	regexp.MustCompile(`.*u[^aeıioöuüâûî]*$`):     `u`,
	regexp.MustCompile(`.*û[lk][^aeıioöuüâûî]*$`): `ü`,
	regexp.MustCompile(`.*û[^lkaeıioöuüâûî]*$`):   `u`,
}

// MakeRoot builds a Root from Latin and Visenc spelling of a word by automatically filling other information
func MakeRoot(latin string, visenc string, pos dervaze.PartOfSpeech) dervaze.Root {

	ow := dervaze.OttomanWord{
		Visenc:           visenc,
		Unicode:          VisencToUnicode(visenc),
		Abjad:            VisencToAbjad(visenc),
		VisencLetters:    VisencLetters(visenc),
		SearchKey:        VisencToSearchKey(visenc),
		DotlessSearchKey: VisencToDotlessSearchKey(visenc),
	}

	r := dervaze.Root{
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

	UpdateEffectiveSoftening(r)

	return r
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
	} else {
		return false
	}
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
func UpdateEffectiveSoftening(r *dervaze.Root) {

	if strings.HasSuffix(r.TurkishLatin, "k") &&
		strings.HasSuffix(r.Ottoman.Visenc, "fo2") {
		tll := len(r.TurkishLatin)
		r.EffectiveTurkishLatin = r.TurkishLatin[0:tll-1] + "ğ"
		ovl := len(r.Ottoman.Visenc)
		r.EffectiveVisenc = r.Ottoman.Visenc[0:ovl-3] + "ao1"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "p") && strings.HasSuffix(r.Ottoman.Visenc, "bu1") {
		tll := len(r.TurkishLatin)
		r.EffectiveTurkishLatin = r.TurkishLatin[0:tll-1] + "b"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "ç") && strings.HasSuffix(r.Ottoman.Visenc, "xu1") {
		tll := len(r.TurkishLatin)
		r.EffectiveTurkishLatin = r.TurkishLatin[0:tll-1] + "c"
		r.HasConsonantSoftening = true
	}

	if strings.HasSuffix(r.TurkishLatin, "t") && strings.HasSuffix(r.Ottoman.Visenc, "d") {

		tll := len(r.TurkishLatin)
		r.EffectiveTurkishLatin = r.TurkishLatin[0:tll-1] + "d"
		r.HasConsonantSoftening = true
	}

}
