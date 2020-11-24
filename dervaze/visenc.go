package dervaze

import "golang.org/x/text/unicode/norm"

var visencToOttoman = map[string]string{
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


func MakeOttomanWord(visenc string, unicode string) (OttomanWord, error) {
	if (visenc == "" && unicode == "") {
		return nil, errors.New("Need either visenc or ottoman")
	}

	if len(unicode) > 0 {
		normalized := norm.NFKC.String(unicode)
	} else {
		normalized := norm.NFKC.String(VisencToUnicode(visenc))
	}

	if len(visenc) > 0 {
		visenc = UnicodeToVisenc(unicode)
	}

	abjad := UnicodeToAbjad(normalized)

	searchKey = SearchKey(visenc)


}

func SearchKey(s string) {
	sk := regexp.MustCompile(`\([oui][0456789]+\)`)
	return sk.ReplaceAllLiteralString(s, "")
}

func DotlessSearchKey(s string) {
	sk := regexp.MustCompile(`\([oui][0123456789]+\)`)
	return sk.ReplaceAllLiteralString(s, "")
}

func SplitVisenc(s string) {

}


OttomanWord getOttomanWord({String visenc, String ottoman}) {
  ow.visencLetters.addAll(splitVisenc(ow.searchKey));
  return ow;
}

func VisencToUnicode(s string) {
  
}

String toOttoman(String v,
    {bool includeInvalid = true, int maxVisencLength = 5}) {
  var vlen = v.length;
  var mx = maxVisencLength;
  var start = 0;
  var end = start + mx < vlen ? start + mx : vlen;

  var res = "";

  while (start < vlen) {
    var possibleLetter = v.substring(start, end);
    if (VISENC_TO_OTTOMAN.containsKey(possibleLetter)) {
      res += VISENC_TO_OTTOMAN[possibleLetter];
      start = end;
      end = start + mx < vlen ? start + mx : vlen;
    } else {
      if (end > (start + 1)) {
        end -= 1;
      } else {
        if (includeInvalid) {
          res += possibleLetter;
        }
        start = start + 1;
        end = start + mx < vlen ? start + mx : vlen;
      }
    }
  }
  return res;
}

String toVisenc(String s,
    {bool includeInvalid = true, int maxSubstringLength = 2}) {
  var slen = s.length;
  var mx = maxSubstringLength;
  var start = 0;
  var end = start + mx < slen ? start + mx : slen;

  var res = "";

  while (start < slen) {
    var possibleLetter = s.substring(start, end);
    if (OTTOMAN_TO_VISENC.containsKey(possibleLetter)) {
      res += OTTOMAN_TO_VISENC[possibleLetter];
      start = end;
      end = start + mx < slen ? start + mx : slen;
    } else {
      if (end > (start + 1)) {
        end -= 1;
      } else {
        if (includeInvalid) {
          res += possibleLetter;
        }
        start = start + 1;
        end = start + mx < slen ? start + mx : slen;
      }
    }
  }
  return res;
}

/// Converts an Arabic/Ottoman string to its abjad representation
int toAbjad(String s) {
  var visenc = toVisenc(s);

  var start = 0;
  var end = start + 3 < visenc.length ? start + 3 : visenc.length;
  var abjad_sum = 0;

  while (start < visenc.length) {
    var possible_letter = visenc.substring(start, end);
    if (VISENC_TO_ABJAD.containsKey(possible_letter)) {
      abjad_sum += VISENC_TO_ABJAD[possible_letter];
      start = end;
      end = start + 3 < visenc.length ? start + 3 : visenc.length;
    } else {
      if (end > start) {
        end -= 1;
      } else {
        start = start + 1;
        end = start + 3 < visenc.length ? start + 3 : visenc.length;
      }
    }
  }

  return abjad_sum;
}

func TFstring(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue } else {
			return ifFalse
		}
}
func TFint(condition bool, ifTrue, ifFalse int) int {
	if condition {
		return ifTrue } else {
			return ifFalse
		}
}

func SplitVisenc(s string, addInvalidLetters bool) []string {
	slen := len(s)
	maxVisencLen := 5
	start := 0
	end := TFint(start + maxVisencLen < slen, start+maxVisencLen + 1, slen + 1)
	group := make([]string, 0, len(s))

	for start < slen {
		visenc, exists := visencToOttoman[s[start:end]]
		if exists {
			group.add(visenc)
			start = end
			end = TFint(start + maxVisencLen < slen, start+maxVisencLen + 1, slen + 1)
		} else {
			if end > start + 1 {
				end -= 1
			} else {
				if addInvalidLetters {
					group.add(s[start:end])
				}
				start += 1
				end = TFint(start + maxVisencLen < slen, start+maxVisencLen + 1, slen + 1)
			}
		}
		
	}

	return group
}

List<String> splitVisenc(String v,
    {bool addInvalidLetters = true, int maxVisencLength = 5}) {
  var vlen = v.length;
  var mx = maxVisencLength;
  var start = 0;
  var end = start + mx < vlen ? start + mx : vlen;

  var elements = <String>[];

  while (start < vlen) {
    var possible_letter = v.substring(start, end);
    if (VISENC_TO_OTTOMAN.containsKey(possible_letter)) {
      elements.add(possible_letter);
      start = end;
      end = start + mx < vlen ? start + mx : vlen;
    } else {
      if (end > start) {
        end -= 1;
      } else {
        if (addInvalidLetters) {
          elements.add(possible_letter);
        }
        start = start + 1;
        end = start + mx < vlen ? start + mx : vlen;
      }
    }
  }

  return elements;
}

var VISENC_TO_OTTOMAN = map[string]string {
  "c": "ء",
  "eo6": "آ",
  // "A": "آ",
  "e": "ا",
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
  "x": "ح",
  "xo1": "خ",
  // "X": "خ",
  "do1": "ذ",
  "d": "د",
  "ro1": "ز",
  "r": "ر",
  // "Z": "ز",
  "ro3": "ژ",
  // "J": "ژ",
  "s": "س",
  "so3": "ش",
  // "S": "ش",
  // "Ş": "ش",
  "z": "ص",
  "zo1": "ض",
  // "D": "ض",
  "t": "ط",
  "to1": "ظ",
  "a": "ع",
  "ao1": "غ",
  // "G": "غ",
  "fo1": "ف",
  // "F": "ف",
  "fo2": "ق",
  // "Q": "ق",
  "lo5": "ك",
  "ko5": "ك",
  "k": "ک",
  "ko7": "گ",
  // "K": "گ",
  "ko3": "ڭ",
  // "lo5o3": "ڭ",
  "l": "ل",
  "m": "م",
  "bo1": "ن",
  // "N": "ن",
  "w": "و",
  "wo5": "ؤ",
  "h": "ه",
  // "h": "ە",
  // "h": "\u06D5",
  "ho2": "ة",
  "y": "ی", // x6cc
  // "y": "ى", // x649
  "bu2": "ي",
  // "Y": "ي",
  "yo5": "ئ",
  "bo5": "ئ",
  "n0": "۰",
  "n1": "۱",
  "n2": "۲",
  "n3": "۳",
  "n4": "۴",
  "n5": "۵",
  "n6": "۶",
  "n7": "۷",
  "n8": "۸",
  "n9": "۹",
  "&zwnj;": "\u200C",
  "||": "\u200C",
  "<>": "\u200C",
  "&zwj;": "\u200D",
  "><": "\u200D",
  "&lrm;": "\u200E",
  "&rlm;": "\u200F",
  "&ls;": "\u2028",
  "&ps;": "\u2028",
  "&lre;": "\u202A",
  "&rle;": "\u202B",
  "&pdf;": "\u202C",
  "&lro;": "\u202D",
  "&rlo;": "\u202D",
  "&bom;": "\uFEFF",
  "o4": "\u064E",
  "u4": "\u0650",
  "o9": "\u064F",
  "u44": "\u064D",
  "o44": "\u064B",
  "o99": "\u064C",
  "o8": "\u0651",
  "o0": "\u0652",
  "o6": "\u0653",
  " ": " ",
  "bot": "\u0679",
  "o5": "\u0654",
  "u5": "\u0655",
}

List<String> visencKeyOrder = VISENC_TO_OTTOMAN.keys.toList()
  ..sort((a, b) => b.length.compareTo(a.length));

String toUnicode(String visenc) {
  var vc = visenc;
  // Log.log('visenc: ' + visenc);
  visencKeyOrder.forEach((v) {
    vc = vc.replaceAll(v, VISENC_TO_OTTOMAN[v]);
  });
  // Log.log('vc: ' + vc);
  return unorm.nfkc(vc);
}

// Checks whether the string contains unicode runes between 0600 and 06FF
bool containsArabicChars(String s) {
  var contains = false;
  s.runes.forEach((u) {
    if (u >= 0x0600 && u <= 0x06FF) {
      contains = true;
    }
  });
  return contains;
}

bool containsDigits(String s) {
  var contains = false;
  s.runes.forEach((u) {
    if (u >= 0x0030 && u <= 0x0039) {
      contains = true;
    }
  });
  return contains;
}

const OTTOMAN_TO_VISENC = {
  'ء': 'c',
  'آ': 'eo6',
  // 'آ': 'A',
  'ا': 'e',
  'أ': 'eo5',
  'إ': 'eu5',
  // 'أ': 'E',
  'ب': 'bu1',
  // 'ب': 'B',
  'پ': 'bu3',
  // 'پ': 'P',
  'ت': 'bo2',
  // 'ت': 'T',
  'ث': 'bo3',
  'ج': 'xu1',
  // 'ج': 'C',
  'چ': 'xu3',
  // 'چ': 'Ç',
  'ح': 'x',
  'خ': 'xo1',
  // 'خ': 'X',
  'د': 'd',
  'ذ': 'do1',
  'ر': 'r',
  'ز': 'ro1',
  // 'ز': 'Z',
  'ژ': 'ro3',
  // 'ژ': 'J',
  'س': 's',
  'ش': 'so3',
  // 'ش': 'S',
  // 'ش': 'Ş',
  'ص': 'z',
  'ض': 'zo1',
  // 'ض': 'D',
  'ط': 't',
  'ظ': 'to1',
  'ع': 'a',
  'غ': 'ao1',
  // 'غ': 'G',
  'ف': 'fo1',
  // 'ف': 'F',
  'ق': 'fo2',
  // 'ق': 'Q',
  // 'ك': 'lo5',
  'ك': 'k',
  'ک': 'k',
  'گ': 'ko7',
  // 'گ': 'K',
  'ڭ': 'ko3',
  // 'ڭ': 'lo5o3',
  'ل': 'l',
  'م': 'm',
  'ن': 'bo1',
  // 'ن': 'N',
  'و': 'w',
  'ؤ': 'wo5',
  'ه': 'h',
  'ە': 'h',
  'ة': 'ho2',
  'ی': 'y', // x6cc
  'ى': 'y', // x649
  'ي': 'bu2',
  // 'ي': 'Y',
  // 'ئ': 'yo5',
  'ئ': 'bo5',
  '۰': 'n0',
  '۱': 'n1',
  '۲': 'n2',
  '۳': 'n3',
  '۴': 'n4',
  '۵': 'n5',
  '۶': 'n6',
  '۷': 'n7',
  '۸': 'n8',
  '۹': 'n9',
  // '\u200C': '&zwnj;',
  '\u200C': '||',
  // '\u200C': '<>',
  // '\u200D': '&zwj;',
  '\u200D': '><',
  '\u200E': '&lrm;',
  '\u200F': '&rlm;',
  '\u2028': '&ls;',
  // '\u2028': '&ps;',
  '\u202A': '&lre;',
  '\u202B': '&rle;',
  '\u202C': '&pdf;',
  '\u202D': '&lro;',
  // '\u202D': '&rlo;',
  '\uFEFF': '&bom;',
  '\u064E': 'o4',
  '\u0650': 'u4',
  '\u064F': 'o9',
  '\u064D': 'u44',
  '\u064B': 'o44',
  '\u064C': 'o99',
  '\u0651': 'o8',
  '\u0652': 'o0',
  '\u0653': 'o6',
  ' ': ' ',
  '\u0679': 'bot',
  '\u0654': 'o5',
  '\u0655': 'u5',
};

List<String> ottomanKeyOrder = OTTOMAN_TO_VISENC.keys.toList()
  ..sort((a, b) => b.length.compareTo(a.length));

String unicodeToVisenc(String ottoman) {
  // Log.log('ottoman: ' + ottoman);
  var ot = unorm.nfkc(ottoman);
  ot = ottoman;
  ottomanKeyOrder.forEach((o) {
    ot = ot.replaceAll(o, OTTOMAN_TO_VISENC[o]);
  });
  // Log.log('ot: ' + ot);
  // Log.log('ot.runes: ' + ot.runes.toString());
  return ot;
}

const VISENC_TO_ABJAD = {
  'e': 1,
  'bu1': 2,
  'bu3': 2,
  'xu1': 3,
  'xu3': 3,
  'd': 4,
  'h': 5,
  'w': 6,
  'ro1': 7,
  'ro3': 7,
  'x': 8,
  't': 9,
  'y': 10,
  'bu2': 10,
  'k': 20,
  'ko7': 20,
  'l': 30,
  'm': 40,
  'bo1': 50,
  's': 60,
  'a': 70,
  'fo1': 80,
  'z': 90,
  'fo2': 100,
  'r': 200,
  'so3': 300,
  'bo2': 400,
  'bo3': 500,
  'xo1': 600,
  'do1': 700,
  'zo1': 800,
  'to1': 900,
  'ao1': 1000,
};