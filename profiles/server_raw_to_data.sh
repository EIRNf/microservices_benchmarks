#!/bin/bash

INPUT_FILE=$1;
TABLE_NAME=$2;

# Figure out table name if table name not provided


#CONVERT TO FOLDED TEXT
pprofutils folded "$INPUT_FILE" $TABLE_NAME.out;

#CONVERT TO CSV
./folded_to_csv.awk $TABLE_NAME.out > $TABLE_NAME.csv;

#CREATE TABLE 
sqlite3 profiles.db "DROP TABLE IF EXISTS $TABLE_NAME;"
#sqlite3 profiles.db "create table $TABLE_NAME";
sqlite3 profiles.db ".import $TABLE_NAME.csv $TABLE_NAME --csv";

#REMOVE INTERMEDIATE FILES  
rm $TABLE_NAME.out $TABLE_NAME.csv

#RUN QUERIES
./process_queries.sh $TABLE_NAME


