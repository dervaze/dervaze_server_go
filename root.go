package root

import (
	"dervazepb"
	"regexp"
	"strings"
)

const VOWELS = "aâeıioöuûü"

var ultimateVowelRegex regexp.Regexp = regexp.MustCompile(`.*([aâeıiîoöuüû])[^aâeıiîoöuüû]*$`)
var ultimateConsonantRegex regexp.Regexp = regexp.MustCompile(`.*([^aâeıiîoöuüû])[aâeıiîoöuüû]*$`)
var vowelsRegex regexp.Regexp = regexp.MustCompile(`([aâeıiîoöuüû])`)
var consonantsRegex regexp.Regexp = regexp.MustCompile(`([bcçdfgğhjklmnpqrsştvwyxz])`)
var endsWithVowelRegex regexp.Regexp = regexp.MustCompile(`.*[aeıioöuüâûî][']?$`)
var hasSingleVowelRegex regexp.Regexp = regexp.MustCompile(`^[^aâeıiîoöuüû]*[aâeıiîoöuüû][^aâeıiîoöuüû]*$`)
var lastConsonantHardRegex regexp.Regexp = regexp.MustCompile(`[fstkçşhp]'?$`)
var effectiveLastVowelRegexes = map[regexp.Regexp]string{
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

func MakeRoot(latin string, visenc string, pos dervazepb.PartOfSpeech) dervazepb.Root {

	ow := OttomanWord{
		Visenc:           visenc,
		Unicode:          VisencToUnicode(visenc),
		Abjad:            VisencToAbjad(visenc),
		VisencLetters:    VisencLetters(visenc),
		SearchKey:        VisencToSearchKey(visenc),
		DotlessSearchKey: VisencToDotlessSearchKey(visenc),
	}

	r := dervazepb.Root{
		TurkishLatin:          latin,
		Ottoman:               ow,
		LastVowel:             LatinLastVowel(latin),
		LastConsonant:         LatinLastConsonant(latin),
		EffectiveLastVowel:    LatinEffectiveLastVowel(latin),
		EffectiveTurkishLatin: LatinEffectiveTurkishLatin(latin),
		EffectiveVisenc:       LatinEffectiveVisenc,
		Abjad:                 ow.Abjad,
		PartOfSpeech:          pos,
		EndsWithVowel:         LatinEndsWithVowel(latin),
		HasSingleVowel:        LatinHasSingleVowel(latin),
		LastVowelHard:         LatinLastVowelHard(latin),
		LastConsonantHard:     LatinLastConsonantHard(latin),
		HasConsonantSoftening: LatinHasConsonantSoftening(latin),
	}
}

func EndsWithVowel(s string) bool {
	return endsWithVowelRegex.MatchString(s)
}

func HasSingleVowel(s string) bool {
	return hasSingleVowelRegex.MatchString(s)
}

func LastConsonantHard(s string) bool {
	return lastConsonantHardRegex.MatchString(s)
}

func LastVowelHard(s string) bool {
	ev := EffectiveLastVowel(s)
	if ev == 'a' || ev == 'ı' || ev == 'o' || ev == 'u' {
		return true
	} else {
		return false
	}
}

func EffectiveLastVowel(s string) string {
	for r, v := range effectiveLastVowelRegexes {
		if r.MatchString(s) {
			return v
		}
	}
}

func LastVowel(s string) string {
	return ultimateVowelRegex.FindString(s)
}
func LastConsonant(s string) string {
	return ultimateConsonantRegex.FindString(s)
}

func (Root *r) UpdateEffectiveSoftening() {
	// these may be modified below according to conditions
	if !(r.hasEffectiveTurkishLatin() &&
		r.hasEffectiveVisenc() &&
		r.hasHasConsonantSoftening()) {
		r.EffectiveTurkishLatin = r.TurkishLatin
		r.EffectiveVisenc = r.Ottoman.Visenc
		r.HasConsonantSoftening = false
	}

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

func LatinEndsWithVowel(latin string) bool {

}
