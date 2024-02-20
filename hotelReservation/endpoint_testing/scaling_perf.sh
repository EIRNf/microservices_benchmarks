#!/bin/bash

target=notnets
experiment_name=scaling_fixed_12-22
# Set the parameters you want to iterate over
num_messages=(1000)
num_instances=(8 16 32 64)

# Set the directory where you want to store the results
runs_directory="runs"

# Create the output directory if it doesn't exist
mkdir -p "$runs_directory"

# Iterate over parameter combinations
for param1 in "${num_messages[@]}"; do
    for param2 in "${num_instances[@]}"; do
        # Generate a unique directory name based on parameters
        run_directory="${runs_directory}/${experiment_name}/${target}_${param1}_${param2}"

        # Create the run directory if it doesn't exist
        mkdir -p "$run_directory"

        # Run your experiment script with parameters and save output to a file
        ./perf_runner.sh "$param1" "$param2" > "${run_directory}/summary.txt" 2>&1

        # Capture output files and move them to directory.
        mv output_* $run_directory/
    done
done
