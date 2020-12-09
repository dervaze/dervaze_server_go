package lang

import (
	context "context"
	"fmt"
	"regexp"
	"strconv"
)

// DervazeServerImpl implementation
type DervazeServerImpl struct {
}

// NewDervazeServerImpl builds a new server instance
func NewDervazeServerImpl() *DervazeServerImpl {
	return &DervazeServerImpl{}
}

func (DervazeServerImpl) mustEmbedUnimplementedDervazeServer() {}

// VisencToOttoman converts a visenc string to Ottoman unicode
func (DervazeServerImpl) VisencToOttoman(ctx context.Context, in *OttomanWord) (*OttomanWord, error) {

	out, err := MakeOttomanWord(in.Visenc, "")
	return out, err

}

// OttomanToVisenc converts a Unicode string to Visenc
func (DervazeServerImpl) OttomanToVisenc(ctx context.Context, in *OttomanWord) (*OttomanWord, error) {
	out, err := MakeOttomanWord("", in.Unicode)
	return out, err
}

// SearchRoots makes a search with various fields and types and returns a Rootset described by the result
func (DervazeServerImpl) SearchRoots(ctx context.Context, in *SearchRequest) (*RootSet, error) {

	var rootList []*Root
	var err error = nil

	searchField := in.SearchField
	searchString := in.SearchString
	maxLen := int(in.ResultLimit)

	switch in.SearchType {
	case SearchType_FUZZY:
		switch searchField {
		case SearchField_AUTO:
			rootList = FuzzySearchAuto(searchString, maxLen)
		case SearchField_OTTOMAN:
			rootList = FuzzySearchUnicode(searchString, maxLen)
		case SearchField_TURKISH_LATIN:
			rootList = FuzzySearchTurkishLatin(searchString, maxLen)
		case SearchField_VISENC:
			rootList = FuzzySearchVisenc(searchString, maxLen)
		case SearchField_ABJAD:
			if s, e := strconv.Atoi(searchString); e == nil {
				rootList = IndexSearchAbjad(int32(s), maxLen)
			} else {
				err = e
			}
		}
	case SearchType_REGEX:
		if searchRegex, e := regexp.Compile(searchString); e == nil {

			switch searchField {
			case SearchField_AUTO:
				rootList = RegexSearchAuto(searchRegex, maxLen)
			case SearchField_OTTOMAN:
				rootList = RegexSearchUnicode(searchRegex, maxLen)
			case SearchField_TURKISH_LATIN:
				rootList = RegexSearchTurkishLatin(searchRegex, maxLen)
			case SearchField_VISENC:
				rootList = RegexSearchVisenc(searchRegex, maxLen)
			case SearchField_ABJAD:
				if s, e := strconv.Atoi(searchString); e == nil {
					rootList = IndexSearchAbjad(int32(s), maxLen)
				} else {
					err = e
				}
			}
		} else {
			err = e
		}

	case SearchType_PREFIX:

		switch searchField {
		case SearchField_AUTO:
			rootList = PrefixSearchAuto(searchString, maxLen)
		case SearchField_OTTOMAN:
			rootList = PrefixSearchUnicode(searchString, maxLen)
		case SearchField_TURKISH_LATIN:
			rootList = PrefixSearchTurkishLatin(searchString, maxLen)
		case SearchField_VISENC:
			rootList = PrefixSearchVisenc(searchString, maxLen)
		case SearchField_ABJAD:
			if s, e := strconv.Atoi(searchString); e == nil {
				rootList = IndexSearchAbjad(int32(s), maxLen)
			} else {
				err = e
			}
		}

	}

	rs := RootSet{Roots: rootList}
	return &rs, err

}

// Translate returns the translation of an Ottoman or Turkish latin sentence
func (DervazeServerImpl) Translate(ctx context.Context, in *TranslateRequest) (*TranslateResponse, error) {
	return nil, fmt.Errorf("Not Implemented")
}
