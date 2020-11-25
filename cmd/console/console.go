package main

import (
	dervaze "dervaze/lang"
	"flag"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// TODO Write real completion with word lists etc
var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	// readline.PcItem("login"),
	// readline.PcItem("say",
	// 	readline.PcItemDynamic(listFiles("./"),
	// 		readline.PcItem("with",
	// 			readline.PcItem("following"),
	// 			readline.PcItem("items"),
	// 		),
	// 	),
	// 	readline.PcItem("hello"),
	// 	readline.PcItem("bye"),
	// ),
	// readline.PcItem("setprompt"),
	// readline.PcItem("setpassword"),
	// readline.PcItem("bye"),
	// readline.PcItem("help"),
	// readline.PcItem("go",
	// 	readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
	// 	readline.PcItem("install",
	// 		readline.PcItem("-v"),
	// 		readline.PcItem("-vv"),
	// 		readline.PcItem("-vvv"),
	// 	),
	// 	readline.PcItem("test"),
	// ),
	// readline.PcItem("sleep"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func console() {

	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		HistoryFile:         "/tmp/dervaze.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "v2o "):
			println(dervaze.VisencToUnicode(line[4:]))
		case strings.HasPrefix(line, "o2v "):
			println(dervaze.UnicodeToVisenc(line[4:]))
		case strings.HasPrefix(line, "t "):
			println(dervaze.PrintRoots(dervaze.SearchTurkishLatin(line[2:])))
		case strings.HasPrefix(line, "v "):
			println(dervaze.PrintRoots(dervaze.SearchVisenc(line[2:])))
		case strings.HasPrefix(line, "u "):
			println(dervaze.PrintRoots(dervaze.SearchUnicode(line[2:])))
		case strings.HasPrefix(line, "a "):
			n, err := strconv.Atoi(line[2:])
			if err != nil {
				println("Need a number for abjad search a ")
			} else {
				println(dervaze.PrintRoots(dervaze.SearchAbjad(int32(n))))
			}
		case line == "":
		default:
			println(dervaze.PrintRoots(dervaze.SearchAll(line)))
		}
	}
}

func main() {

	var inputfile string
	flag.StringVar(&inputfile, "i", "assets/dervaze-rootset.protobuf", "protobuffer file to load roots")

	flag.Parse()
	dervaze.InitSearch(inputfile)
	console()

}
