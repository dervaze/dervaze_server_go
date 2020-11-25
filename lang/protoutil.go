package lang


const PROTOBUFFILENAME := "assets/dervaze-rootset.protobuf"

func SaveRootSetProtobuf(rootset RootSet) {
	// t := time.Now().Format("2006-01-02-03-04-05")
	// filename := fmt.Sprintf("dervaze-roots-%s.bin", t)
	filename := PROTOBUFFILENAME
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteSlice, err := proto.Marshal(&rootset)

	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err := file.Write(byteSlice)

		if err != nil {
			log.Fatal(err)
		}


	log.Printf("%s: Wrote %d bytes.\n", filename, bytesWritten)

	// Write bytes to file
	// totalBytes := 0
	// for i, r := range roots {
	// 	byteSlice, err := proto.Marshal(&r)
	// 	if err != nil {
	// 		// log.Println(string(r.Ottoman.Unicode))
	// 		// log.Printf("Ottoman.Unicode: %x %t ", r.Ottoman.Unicode, utf8.ValidString(r.Ottoman.Unicode))
	// 		// log.Printf("Ottoman.Visenc: %x %t", r.Ottoman.Visenc, utf8.ValidString(r.Ottoman.Visenc))
	// 		// log.Printf("Turkish Latin %x %t", r.TurkishLatin, utf8.ValidString(r.TurkishLatin))
	// 		log.Println(r)
	// 		// log.Println(utf8.ValidString(r.TurkishLatin))
	// 		// log.Println(utf8.ValidString(strings.Join(r.Ottoman.VisencLetters, "")))
	// 		// log.Println(utf8.ValidString(r.Ottoman.SearchKey))
	// 		// log.Println(utf8.ValidString(r.Ottoman.DotlessSearchKey))
	// 		// log.Println(utf8.ValidString(r.LastVowel))
	// 		// log.Println(utf8.ValidString(r.LastConsonant))
	// 		// log.Println(utf8.ValidString(r.EffectiveLastVowel))
	// 		log.Println(r.EffectiveTurkishLatin)
	// 		log.Println(utf8.ValidString(r.EffectiveTurkishLatin))
	// 		// log.Println(utf8.ValidString(r.EffectiveVisenc))
	// 		// log.Println(proto.MarshalTextString(&r))
	// 		log.Println(err)
	// 	}
	// 	bytesWritten, err := file.Write(byteSlice)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	totalBytes += bytesWritten
	// 	if i%1000 == 0 {
	// 		fmt.Printf("%d\n", i)
	// 	}
	// }
	// log.Printf("%s: Wrote %d bytes.\n", filename, totalBytes)
}


func LoadRootSetProtobuf() RootSet {
}
