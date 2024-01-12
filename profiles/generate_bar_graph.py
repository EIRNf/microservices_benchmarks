import csv
import sys
import os
import matplotlib.pyplot as plt
from matplotlib.colors import ListedColormap
import numpy as np

__location__ = os.path.realpath(
    os.path.join(os.getcwd(), os.path.dirname(__file__)))

# Check if the CSV file name is provided as a command line argument
if len(sys.argv) < 2:
    print("Provide the CSV file name as a command line argument.")
    sys.exit(1)

# Get the CSV file name from the command line argument
csv_file = sys.argv[1]
name = csv_file.removesuffix(".csv")

# Read the CSV file
data = []
with open(os.path.join(__location__,csv_file), 'r') as file:
    reader = csv.reader(file)
    for row in reader:
        data.append(row)

labels = []
values = []

# Separate labels and values
for row in data:
    if (row[0] == name):
        continue
    if (row[0] == "Total Samples"):
        continue
    if (row[0] == "Total gRPC"):
        continue
    labels.append(row[0])
    values.append(int(float(row[1]))) 

# Calculate the total sum of values
# for 
total = sum(values)

# Calculate the percentage values
percentages = [(value / total) * 100 for value in values]


# Create a colormap with different colors for each category
colors = ListedColormap(["#58b5e1", "#1f4196", "#8bda64", "#0a4f4e", "#56ebd3"])

# Create a bar graph with different colors for each bar
plt.bar(labels, percentages, color=colors(np.arange(len(labels))))


plt.xlabel('Categories')
plt.ylabel('% Total')
plt.title(name)

# Rotate the x-axis labels for better readability
plt.xticks(rotation=15)

# Set the y-axis tick labels as percentages
plt.yticks(np.arange(0, max(percentages)+10, 10))
plt.gca().set_yticklabels(['{:.1f}%'.format(val) for val in plt.gca().get_yticks()])

# Add light grid lines
plt.grid(axis='y', alpha=0.3)

# Save the graph as an SVG file
plt.savefig(name + '.svg', format='svg')

# Display the graph
# plt.show()