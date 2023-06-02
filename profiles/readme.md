


To transform the pprof data to csv file, inport into a sqlite db for queries, and generate resulting graphs run the following.

```
./process_folder.sh raw_data
```


To simply regenerate a particular graph, run the generate_bar_graph.py script with a csv file as input. eg.

```
python3 generate_bar_graph.py csv_data/profile.csv 
```

## Issues

- The frontend and search services are not accurately analysed. 
    - As they function as GRPC clients and servers it is possible that the server queries capture Client functions as well rending analysis innacurate.
