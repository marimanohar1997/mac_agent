package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type SystemInfo struct {
	Hostname        string           `json:"hostname"`
	OS              string           `json:"os"`
	Platform        string           `json:"platform"`
	PlatformVersion string           `json:"platform_version"`
	CPUModel        string           `json:"cpu_model"`
	CPUCores        int              `json:"cpu_cores"`
	TotalMemory     uint64           `json:"total_memory"`
	UsedMemory      uint64           `json:"used_memory"`
	TotalDisk       uint64           `json:"total_disk"`
	UsedDisk        uint64           `json:"used_disk"`
	InstalledApps   []AppInfo        `json:"installed_apps"`
	OpenPorts       []PortInfo       `json:"open_ports"`
	ActiveConns     []ConnectionInfo `json:"active_connections"`
	Timestamp       time.Time        `json:"timestamp"`
}

type AppInfo struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	InstallDate  time.Time `json:"install_date"`
	ObtainedFrom string    `json:"obtained_from"`
}

type PortInfo struct {
	Port     uint32 `json:"port"`
	Protocol string `json:"protocol"`
}

type ConnectionInfo struct {
	LocalAddress  string `json:"local_address"`
	LocalPort     uint32 `json:"local_port"`
	RemoteAddress string `json:"remote_address"`
	RemotePort    uint32 `json:"remote_port"`
	Status        string `json:"status"`
	Type          string `json:"type"`
}

func main() {
	info := collectSystemInfo()
	saveToFile(info)
	sendToAPI(info)
}

func collectSystemInfo() SystemInfo {
	hostInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	return SystemInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		CPUModel:        cpuInfo[0].ModelName,
		CPUCores:        runtime.NumCPU(),
		TotalMemory:     memInfo.Total,
		UsedMemory:      memInfo.Used,
		TotalDisk:       diskInfo.Total,
		UsedDisk:        diskInfo.Used,
		InstalledApps:   getInstalledApps(),
		OpenPorts:       getOpenPorts(),
		ActiveConns:     getActiveConnections(),
		Timestamp:       time.Now(),
	}
}

func getInstalledApps() []AppInfo {
	cmd := exec.Command("system_profiler", "SPApplicationsDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting installed apps: %v", err)
		return []AppInfo{}
	}

	var result map[string]interface{}
	err = json.Unmarshal(output, &result)
	if err != nil {
		log.Printf("Error parsing installed apps JSON: %v", err)
		return []AppInfo{}
	}

	appsData, ok := result["SPApplicationsDataType"].([]interface{})
	if !ok {
		log.Printf("Error: SPApplicationsDataType is not a slice")
		return []AppInfo{}
	}

	var installedApps []AppInfo
	for _, app := range appsData {
		appInfo, ok := app.(map[string]interface{})
		if !ok {
			log.Printf("Error: app is not a map")
			continue
		}

		name, _ := appInfo["_name"].(string)
		version, _ := appInfo["version"].(string)

		var installDate time.Time
		if lastModifiedStr, ok := appInfo["lastModified"].(string); ok {
			installDate, err = time.Parse(time.RFC3339, lastModifiedStr)
			if err != nil {
				log.Printf("Error parsing install date for %s: %v", name, err)
			}
		}

		obtainedFrom, _ := appInfo["obtained_from"].(string)

		installedApps = append(installedApps, AppInfo{
			Name:         name,
			Version:      version,
			InstallDate:  installDate,
			ObtainedFrom: obtainedFrom,
		})
	}

	return installedApps
}

func getOpenPorts() []PortInfo {
	connections, err := net.Connections("all")
	if err != nil {
		log.Printf("Error fetching network connections: %v", err)
		return []PortInfo{}
	}

	openPorts := make(map[uint32]string)
	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			openPorts[conn.Laddr.Port] = protocolToString(conn.Type)
		}
	}

	var portInfos []PortInfo
	for port, protocol := range openPorts {
		portInfos = append(portInfos, PortInfo{Port: port, Protocol: protocol})
	}
	return portInfos
}

func getActiveConnections() []ConnectionInfo {
	connections, err := net.Connections("all")
	if err != nil {
		log.Printf("Error fetching network connections: %v", err)
		return []ConnectionInfo{}
	}

	var activeConns []ConnectionInfo
	for _, conn := range connections {
		if conn.Status == "ESTABLISHED" {
			activeConns = append(activeConns, ConnectionInfo{
				LocalAddress:  conn.Laddr.IP,
				LocalPort:     conn.Laddr.Port,
				RemoteAddress: conn.Raddr.IP,
				RemotePort:    conn.Raddr.Port,
				Status:        conn.Status,
				Type:          protocolToString(conn.Type),
			})
		}
	}
	return activeConns
}

func protocolToString(protocol uint32) string {
	switch protocol {
	case 6:
		return "TCP"
	case 17:
		return "UDP"
	default:
		return fmt.Sprintf("Unknown (%d)", protocol)
	}
}

func saveToFile(info SystemInfo) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	outputDir := filepath.Join(homeDir, "Library", "Logs", "mac-agent")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("system_info_%s.json", time.Now().Format("20060102_150405")))

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("System information saved to: %s\n", filename)
}

func sendToAPI(info SystemInfo) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	// Replace this URL with your actual mock API endpoint
	url := "https://mockapi.example.com/system-info"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending data to API: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading API response: %v", err)
		return
	}

	fmt.Printf("API Response: %s\n", string(body))
}
