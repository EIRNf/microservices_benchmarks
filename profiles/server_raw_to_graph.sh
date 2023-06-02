#!/bin/bash

set -eEuo pipefail

INPUT_FILE=$1;
CSV_DIR=$2;
GRAPH_DIR=$3;

delimiter="."
TABLE_NAME="${INPUT_FILE%%$delimiter*}"

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
rm "$INPUT_FILE" $TABLE_NAME.out $TABLE_NAME.csv

#RUN QUERIES
./process_queries.sh $TABLE_NAME > $TABLE_NAME.csv

#GENERATE BAR GRAPH
python3 generate_bar_graph.py $TABLE_NAME.csv

#MOVE CSV OUTPUT FILES TO CHOSEN DIRECTORY
mv $TABLE_NAME.csv $CSV_DIR

mv $TABLE_NAME.svg $GRAPH_DIR




