import os
import json
import numpy as np

# def get_operations(file_path):
#     with open(file_path, 'r') as file:
#         trace_data = json.load(file)
#         operation = span['operationName'] for span in trace_data

#         return calculate_mean_duration(trace_data)


def calculate_mean_duration(trace_data):
    durations = [span['duration'] for span in trace_data['data']]
    return np.mean(durations)

def calculate_total_duration(service,trace_data):
    total_duration = 0
    for span in trace_data['spans']:
        if span['operationName'] == service[0]:
            tag = span['tags'][0]
            if tag['value'] == 'server':
                total_duration += span['duration']
    return total_duration

def calculate_instance_count(service,trace_data):
    instance_count = 0
    for span in trace_data['spans']:
        if span['operationName'] == service[0]:
            tag = span['tags'][0]
            if tag['value'] == 'server':
                instance_count += 1
    return instance_count

def get_durations(service,trace_data):
    durations = []
    for span in trace_data['spans']:
        if span['operationName'] == service[0]:
            tag = span['tags'][0]
            if tag['value'] == 'server':
                duration = span['duration']
                durations.append(duration)
    return durations

def process_trace_file(service_keys,file_path):
    with open(file_path, 'r') as file:
        trace_data = json.load(file)
        for service in service_keys.items():
            duration_list = get_durations(service,trace_data)
            service[1]['instances'].extend(duration_list)
            service[1]['total_duration'] += calculate_total_duration(service, trace_data)
            service[1]['instance_count'] += calculate_instance_count(service, trace_data)

def process_directory(directory_path):
    service_stats = {
        "/reservation.Reservation/CheckAvailability": {
            "instances" : [],
            "instance_count": 0,
            "total_duration": 0,
        },
        "/profile.Profile/GetProfiles": {
            "instances" : [],
            "instance_count": 0,
            "total_duration": 0,
        },
        "/search.Search/Nearby": {
            "instances" : [],
            "instance_count": 0,
            "total_duration": 0,
        },
        "/geo.Geo/Nearby": {
            "instances" : [],
            "instance_count": 0,
            "total_duration": 0,
        },
        "/rate.Rate/GetRates": {
            "instances" : [],
            "instance_count": 0,
            "total_duration": 0,
        },
    }


    mean_durations = {}
    std_dev = {}
    for filename in os.listdir(directory_path):
        if filename.endswith('.json'):
            file_path = os.path.join(directory_path, filename)
            process_trace_file(service_stats,file_path)
            # mean_durations[service_name] = mean_durations[service_Zname] + mean_duration

    for service in service_stats.items():
        # for instance in service_stats[i].keys():
        #     instance['toa']
        mean_durations[service[0]] = service[1]['total_duration']/service[1]['instance_count']
        durations =  service[1]['instances']
        std_dev[service[0]] = np.std(durations)
    return mean_durations, std_dev

def main():
    root_directory = '/home/esiramos/projects/microservices_benchmarks/hotelReservation/endpoint_testing/jeager/frontend' 
    services_mean_durations, std_dev = process_directory(root_directory)

    print(root_directory)
    print("Mean Durations for each service:")
    means = []
    for service, mean_duration in services_mean_durations.items():
        stddev = std_dev[service]
        means.append(mean_duration)
        print(f"{service}: {mean_duration} us, Stddev: {stddev}")

    overall = np.mean(means)
    print(f"Overall Mean: {overall} us")
    # stdev = np.std(a)
    # print(f"Overall Mean: {overall} us")

if __name__ == "__main__":
    main()
