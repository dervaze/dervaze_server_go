#!/bin/zsh


F=${1}
wc -l $F
for h in {a..z} â î û ü ş ç ö ı ğ İ ; do
    rg -i "^${h}" $F | sort > ${F:r}-${h}.csv 
done

rg -i "^[^a-zâîûüşöçığİ]" $F | sort > ${F:r}-noletter.csv

cat ${F:r}-â.csv >> ${F:r}-a.csv ; rm -f ${F:r}-â.csv
cat ${F:r}-î.csv >> ${F:r}-i.csv ; rm -f ${F:r}-î.csv
cat ${F:r}-û.csv >> ${F:r}-u.csv ; rm -f ${F:r}-û.csv
cat ${F:r}-İ.csv >> ${F:r}-i.csv ; rm -f ${F:r}-İ.csv


wc -l ${F:r}-*

