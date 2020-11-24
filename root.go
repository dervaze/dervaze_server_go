package dervaze_server

import "dervaze.com/dervazepb"
import "regexp"
import "strings"


const VOWELS = "aâeıioöuûü"
var ultimateVowelRegex := regexp.MustCompile(`.*([aâeıiîoöuüû])[^aâeıiîoöuüû]*$`)
var ultimateConsonantRegex := egexp.MustCompile(`.*([^aâeıiîoöuüû])[aâeıiîoöuüû]*$`)
var vowelsRegex := regexp.MustCompile(`([aâeıiîoöuüû])`)
var consonantsRegex := regexp.MustCompile(`([bcçdfgğhjklmnpqrsştvwyxz])`)
var endsWithVowelRegex := regexp.MustCompile(`.*[aeıioöuüâûî][']?$`)
var hasSingleVowelRegex := regexp.MustCompile(`^[^aâeıiîoöuüû]*[aâeıiîoöuüû][^aâeıiîoöuüû]*$`)
var lastConsonantHardRegex := regexp.MustCompile(`[fstkçşhp]'?$`)
var effectiveLastVowelRegexes := {
  regexp.MustCompile(`.*a[^aeıioöuüâûî]*$`): `a`,
  regexp.MustCompile(`.*â[lk][^aeıioöuüâûî]*$`): `i`,
  regexp.MustCompile(`.*â[^lkaeıioöuüâûî]*$`: `a`,
  regexp.MustCompile(`.*e[^aeıioöuüâûî]*$`): `e`,
  regexp.MustCompile(`.*i[^aeıioöuüâûî]*$`): `i`,
  regexp.MustCompile(`.*î[^aeıioöuüâûî]*$`): `i`,
  regexp.MustCompile(`.*ı[^aeıioöuüâûî]*$`): `ı`,
  regexp.MustCompile(`.*ö[^aeıioöuüâûî]*$`): `ö`,
  regexp.MustCompile(`.*o[^aeıioöuüâûî]*$`): `o`,
  regexp.MustCompile(`.*ü[^aeıioöuüâûî]*$`): `ü`,
  regexp.MustCompile(`.*u[^aeıioöuüâûî]*$`): `u`,
  regexp.MustCompile(`.*û[lk][^aeıioöuüâûî]*$`): `ü`,
  regexp.MustCompile(`.*û[^lkaeıioöuüâûî]*$`): `u`,
}

func MakeRoot(latin string, visenc string, pos dervazepb.PartOfSpeech) dervazepb.Root
{

	ow := OttomanWord{
		Visenc: visenc,
		Unicode: VisencToUnicode(visenc),
		Abjad: VisencToAbjad(visenc),
		VisencLetters: VisencLetters(visenc),
		SearchKey: VisencToSearchKey(visenc),
		DotlessSearchKey: VisencToDotlessSearchKey(visenc)
	}

	r := dervazepb.Root{
		TurkishLatin: latin,
		Ottoman: ow,
		LastVowel: LatinLastVowel(latin),
		LastConsonant: LatinLastConsonant(latin),
		EffectiveLastVowel: LatinEffectiveLastVowel(latin),
		EffectiveTurkishLatin: LatinEffectiveTurkishLatin(latin),
		EffectiveVisenc: LatinEffectiveVisenc,
		Abjad: ow.Abjad,
		PartOfSpeech: pos,
		EndsWithVowel: LatinEndsWithVowel(latin),
		HasSingleVowel: LatinHasSingleVowel(latin),
		LastVowelHard: LatinLastVowelHard(latin),
		LastConsonantHard: LatinLastConsonantHard(latin),
		HasConsonantSoftening: LatinHasConsonantSoftening(latin)
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
  var ev := EffectiveLastVowel(s)
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
	  strings.HasSuffix(r.Ottoman.Visenc, "fo2") 
	  {
      r.effectiveTurkishLatin =
          r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'ğ';
      r.effectiveVisenc =
          r.ottoman.visenc.substring(0, r.ottoman.visenc.length - 3) + 'ao1';
      r.hasConsonantSoftening = true;

	  }
    if (r.turkishLatin.endsWith('k') && r.ottoman.visenc.endsWith('fo2')) {
    }
    if (r.turkishLatin.endsWith('k') && r.ottoman.visenc.endsWith('k')) {
      r.effectiveTurkishLatin =
          r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'ğ';
      r.effectiveVisenc =
          r.ottoman.visenc.substring(0, r.ottoman.visenc.length - 1) + 'ao1';
      r.hasConsonantSoftening = true;
    }
    if (r.turkishLatin.endsWith('p') && r.ottoman.visenc.endsWith('bu1')) {
      r.effectiveTurkishLatin =
          r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'b';
      r.hasConsonantSoftening = true;
    }
    if (r.turkishLatin.endsWith('ç') && r.ottoman.visenc.endsWith('xu1')) {
      r.effectiveTurkishLatin =
          r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'c';
      r.hasConsonantSoftening = true;
      r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'c';
    }
    if (r.turkishLatin.endsWith('t') && r.ottoman.visenc.endsWith('d')) {
      r.effectiveTurkishLatin =
          r.turkishLatin.substring(0, r.turkishLatin.length - 1) + 'd';
      r.hasConsonantSoftening = true;
    }
  }
}



func LatinEndsWithVowel(latin string) bool {

}


