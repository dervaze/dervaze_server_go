#!/usr/bin/python3

import pymongo
from bson import Decimal128, decode, decode_all, decode_iter

import os
import re
import datetime as dt
import csv
from pprint import pprint
import glob

print("Hello")

client = pymongo.MongoClient(
    "mongodb+srv://caribou1:nBljH0U07v7DZQpE@caribou-0.gy3lm.mongodb.net/caribou?retryWrites=true&w=majority"
)

db = client.dervaze

verb_files = glob.glob('v/*.csv')
noun_files = glob.glob('n/*.csv')
proper_files = glob.glob('p/*.csv')

pprint(verb_files)
pprint(noun_files)
pprint(proper_files)

for vf in verb_files[:1]:
    print(f"** {vf} **")
    with open(vf) as vff:
        csvf = csv.reader(vff)
        for row in csvf:
            print(row)


