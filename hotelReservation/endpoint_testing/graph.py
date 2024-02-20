import os
import re
import matplotlib.pyplot as plt
import numpy as np

def parse_output_file(file_path):
    with open(file_path, 'r') as file:
        content = file.read()

        # Define the variables to extract
        variables = ['avg', 'min', 'p50', 'max', 'p90', 'p99', 'p999', 'p9999']

        # Initialize a dictionary to store the values
        data = {variable: 0.0 for variable in variables}

        # Extract values for each variable
        for variable in variables:
            match = re.search(f'{variable}\s+([\d.]+)ms', content)
            if match:
                data[variable] = float(match.group(1))
    return data

def Average(lst): 
    return sum(lst) / len(lst) 

def process_experiment(experiment_directory):
    data_series = {}
    for system_type_dir in os.listdir(experiment_directory):
        system_type_path = os.path.join(experiment_directory, system_type_dir)
        if os.path.isdir(system_type_path):
            system_type_data = process_system_type(system_type_path)
            
            if system_type_dir not in data_series:
                data_series[system_type_dir] = { 'avg': [], 'p50': [], 'p90': [], 'p9999': []}
            
            data_series[system_type_dir]['avg'] = Average(system_type_data['avg'])
            data_series[system_type_dir]['p50'] = Average(system_type_data['p50'])
            data_series[system_type_dir]['p90'] = Average(system_type_data['p90'])
            data_series[system_type_dir]['p9999'] = Average(system_type_data['p9999'])

                
                # data_series[system_type]['avg'].extend(data['avg'])
                # data_series[system_type]['p50'].extend(data['p50'])
                # data_series[system_type]['p90'].extend(data['p90'])
                # data_series[system_type]['p9999'].extend(data['p9999'])

    return data_series

def process_system_type(system_type_directory):
    data = {'avg': [], 'p50': [], 'p90': [], 'p9999': []}
    for filename in os.listdir(system_type_directory):
        if filename.startswith("output_") and filename.endswith(".txt"):
            file_path = os.path.join(system_type_directory, filename)
            file_data = parse_output_file(file_path)

            data['avg'].append(np.mean(file_data.get('avg', [0])))
            data['p50'].append(np.mean(file_data.get('p50', [0])))
            data['p90'].append(np.mean(file_data.get('p90', [0])))
            data['p9999'].append(np.mean(file_data.get('p9999', [0])))

    return data

def plot_line_graph(data_series, experiment_name):

    x_order = ["8","16","32","64"]

    navg = []
    np50 = []
    np90 = []
    np9999 = []
    bavg = []
    bp50 = []
    bp90 = []
    bp9999 = []

    for system_type, data in data_series.items():
        x_axis = system_type.split('_')[2]
        pos = x_order.index(x_axis)
        # x_axis = system_type[len(system_type) - 2]

        if "notnets" in system_type:
           navg.insert(pos, data['avg'])
           np50.insert(pos, data['p50'])
           np90.insert(pos, data['p90'])
           np9999.insert(pos, data['p9999'])

        if "baseline" in system_type:
            bavg.insert(pos, data['avg'])
            bp50.insert(pos, data['p50'])
            bp90.insert(pos, data['p90'])
            bp9999.insert(pos, data['p9999'])
           
    plt.plot(x_order, navg, label=f'notnets-avg', color='blue', linestyle='-')
    plt.plot(x_order, np50, label=f'notnets-p50', color='dodgerblue', linestyle='--')
    plt.plot(x_order, np90, label=f'notnets-p90', color='deepskyblue', linestyle='-.')
    plt.plot(x_order, np9999, label=f'notnets-p9999', color='lightsteelblue', linestyle=':')

    plt.plot(x_order, bavg, label=f'baseline-avg', color='red', linestyle='-')
    plt.plot(x_order, bp50, label=f'baseline-p50', color='tomato', linestyle='--')
    plt.plot(x_order, bp90, label=f'baseline-p90', color='indianred', linestyle='-.')
    plt.plot(x_order, bp9999, label=f'baseline-p9999', color='lightcoral', linestyle=':')
    

    plt.xlabel('Concurrent Clients')
    plt.ylabel('Latency (ms/req)')
    plt.legend(loc='upper left', bbox_to_anchor=(1, 1))
    plt.title('Latency with increasing Clients')
    plt.savefig(experiment_name + ".png",bbox_inches='tight')
    # plt.show()

if __name__ == "__main__":
    experiment_directory = "/home/estebanramos/microservices_benchmarks/hotelReservation/endpoint_testing/runs/scaling_rand"

    experiment_name = os.path.basename(os.path.normpath(experiment_directory))

    data_series = process_experiment(experiment_directory)
    plot_line_graph(data_series, experiment_name)
