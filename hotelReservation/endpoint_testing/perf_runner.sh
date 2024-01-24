#!/bin/bash

# Set the number of instances you want to run
num_instances=1

run_program() {
    instance_id=$1
    input_param=$2

# -benchtime=60s
    go test -test.run=NOP -bench=. $input_param > output_$instance_id.txt
}

awk_cmd(){
    search_string=$1
    file=$2
    awk -F':' -v search="$search_string" '$0 ~ search {gsub(/^[ \t]+/, "", $2); print $2}' "$file"
}

calculate_values() {
    total_requests=0
    longest_execution=0
    sum_throughput=0
    sum_avg_latency=0

    for file in output_*.txt; do
        NUM_REQUESTS=$(awk_cmd "NUM_REQUESTS" "$file")
        EXECUTION_LENGTH=$(awk_cmd "EXECUTION_LENGTH" "$file")
        EXECUTION_LENGTH=$(echo "$EXECUTION_LENGTH" | tr -cd '0-9.')
        THROUGHPUT=$(awk_cmd "THROUGHPUT" "$file")
        THROUGHPUT=$(echo "$THROUGHPUT" | tr -cd '0-9.')
        AVERAGE_LATENCY=$(awk_cmd "AVERAGE_LATENCY" "$file")
        AVERAGE_LATENCY=$(echo "$AVERAGE_LATENCY" | tr -cd '0-9.')

        # Total requests
        total_requests=$(echo print $total_requests + $NUM_REQUESTS | perl)

        # Longest execution
        if perl -e "exit($EXECUTION_LENGTH > $longest_execution? 0 : 1)" ; then
            longest_execution=$EXECUTION_LENGTH
        fi
        # if [ $longest_execution -lt $EXECUTION_LENGTH ]; then
        #     longest_execution=$EXECUTION_LENGTH
        # fi

        # Averaged throughput
        sum_throughput=$(echo print $sum_throughput + $THROUGHPUT | perl)
        # Avg latency
        sum_avg_latency=$(echo print $sum_avg_latency + $AVERAGE_LATENCY | perl)
    done

    echo "Num Instances: $num_instances"
    echo "Total Requests: $total_requests"
    echo "Longest Execution: $longest_execution"
    # echo "Sum Throughput: $sum_throughput"
    echo "Calculated Throughput(longest_execution): $(echo print $total_requests/$longest_execution | perl)"
    echo "Calculated Throughput(runtime): $(echo print $total_requests/$runtime | perl)"

    echo "Averaged Throughput per Instance: $(echo print $sum_throughput/$num_instances | perl)"
    echo "Averaged Latency per Instance: $(echo print $sum_avg_latency/$num_instances | perl)"
}


# Clear outputs 
prefix="output_"
rm -f "$prefix"*

start=`date +%s.%N`
for ((i=1; i<=$num_instances; i++)); do
    run_program $i "-url=http://192.168.49.2:30252 -endpoint=hotels -reqs=6000" &
done
wait
end=`date +%s.%N`

runtime=$( echo print $end - $start | perl )


# Calculate values
calculate_values 




