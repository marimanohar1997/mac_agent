# System Information Collector

## Overview

The System Information Collector is a Go program that gathers detailed information about your computer system, including hardware specifications, installed software, network connections, and more. It saves this information to a JSON file and can also send it to a specified API endpoint.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

1. **Go**: This project is written in Go. If you don't have Go installed, download and install it from [the official Go website](https://golang.org/doc/install).

2. **Git**: You'll need Git to clone the repository. Download and install it from [the official Git website](https://git-scm.com/downloads).

## Installation

1. Open a terminal or command prompt.

2. Clone the repository:

   ```
   git clone https://github.com/marimanohar1997/mac_agent.git
   ```

3. Navigate to the project directory:

   ```
   cd mac_agent
   ```

4. Install the required dependencies:
   ```
   go get github.com/shirou/gopsutil
   ```

## Configuration

1. Open the `main.go` file in a text editor.

2. Locate the `sendToAPI` function (around line 250).

3. Replace the mock API URL with your actual API endpoint:
   ```go
   url := "https://your-actual-api-endpoint.com/mac_agent"
   ```

## Usage

To run the System Information Collector:

1. Open a terminal or command prompt.

2. Navigate to the project directory if you're not already there.

3. Run the program:

   ```
   go run main.go
   ```

4. The program will collect system information and perform two actions:

   - Save the information to a JSON file in your home directory under `Library/Logs/mac_agent/`.
   - Send the information to the configured API endpoint.

5. Check the console output for the location of the saved JSON file and the API response.

## Output

The program generates two types of output:

1. **JSON File**: A detailed JSON file containing all collected system information. You can find this file in your home directory under `Library/Logs/mac_agent/`. The filename includes a timestamp.

2. **API Response**: The program will print the response from the API to the console.

## Troubleshooting

If you encounter any issues:

1. Ensure all prerequisites are correctly installed.
2. Check that you have an active internet connection for sending data to the API.
3. Verify that the API endpoint URL is correct in the `main.go` file.
4. If you get permission errors, try running the program with administrator privileges.

## Contributing

Contributions to improve the System Information Collector are welcome. Please feel free to submit pull requests or create issues for bugs and feature requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

# mac_agent
