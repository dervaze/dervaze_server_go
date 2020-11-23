#!/usr/bin/env python3

import sys
import csv
import json
import codecs
import os

csvfilename = sys.argv[1]
fieldnames = ["turkish_latin", "visenc"]
to_json = []
with codecs.open(csvfilename, "r", "utf-8") as fin:
    csv_reader = csv.DictReader(fin, fieldnames=fieldnames, dialect="unix")
    # skip first line
    csv_reader.__next__()
    for csv_line in csv_reader:
        to_json.append(
            {"turkish_latin": csv_line["turkish_latin"], "visenc": csv_line["visenc"]}
        )

jsonfilename = os.path.splitext(os.path.abspath(csvfilename))[0] + ".json"

with codecs.open(jsonfilename, "w", "utf-8") as fout:
    json.dump(to_json, fout, indent=2)
