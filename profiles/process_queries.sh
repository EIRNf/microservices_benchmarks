#!/bin/bash

set -eEuo pipefail

TABLE_NAME=$1
INCLUDE_TOTALS="false"
# echo "Running queries for table";

# echo $TABLE_NAME;
if [[ $INCLUDE_TOTALS == "true" ]]; then

echo -n "Total Samples,";
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME};";

echo -n "Total gRPC,";
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%grpc%';";
fi


echo -n "gRPC Syscall,";
#sqlite3 profiles.db "CREATE VIEW all_syscall AS SELECT * from ${TABLE_NAME} WHERE CallStack LIKE '%syscall%';";
sqlite3 profiles.db "select SUM(\"Count\") from ${TABLE_NAME} WHERE CallStack LIKE '%syscall%' AND CallStack NOT LIKE '%${TABLE_NAME}.(%';";


echo -n "Business Call,";
if [[ $TABLE_NAME == "merged" ]]; then 
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%\_Handler%' ESCAPE '\';";
else
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%${TABLE_NAME}.(%';";
fi

#Select SUM("Count") from merged where CallStack LIKE '%proto%'AND CallStack NOT LIKE '%client%';

echo -n "Data Transform,";
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%Unmarshal%' OR CallStack Like '%Marshal%' OR CallStack Like '%Decoder%' OR CallStack Like '%Encoder%';";

echo -n "Core,";
#sqlite3 profiles.db "Select SUM(\"Count\") from ( select * from ${TABLE_NAME} where CallStack LIKE '%controlBuffer%' union  select * from ${TABLE_NAME} where CallStack Like '%loopyWriter%' except select * from all_syscall);";
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%controlBuffer%' OR CallStack LIKE '%loopyWriter%' AND CallStack NOT LIKE '%syscall%';"

echo -n "Transport,";
sqlite3 profiles.db "Select SUM(\"Count\") from ${TABLE_NAME} where CallStack LIKE '%http2Server).Write%' OR CallStack LIKE '%http2Server).WriteStatus%' OR CallStack LIKE '%Stream).Read%' OR (CallStack LIKE '%HandleStreams%' AND CallStack NOT LIKE '%syscall%');"
