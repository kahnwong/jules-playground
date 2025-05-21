package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func getSystemTemperature() (float64, error) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0, fmt.Errorf("failed to get sensor temperatures: %w", err)
	}

	if len(temps) == 0 {
		return 0, fmt.Errorf("system temperature not available: no sensors found")
	}

	// Try to find a CPU-related temperature
	for _, temp := range temps {
		// Common keys for CPU temperature sensors
		if strings.Contains(strings.ToLower(temp.SensorKey), "core") || strings.Contains(strings.ToLower(temp.SensorKey), "cpu") || strings.Contains(strings.ToLower(temp.SensorKey), "thermal") {
			return temp.Temperature, nil
		}
	}

	// If no specific CPU sensor is found, return the first available one
	return temps[0].Temperature, nil
}

func getDiskUtilization() (uint64, uint64, float64, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return 0, 0, 0.0, err
	}
	return usage.Total, usage.Used, usage.UsedPercent, nil
}

func getMemoryUtilization() (uint64, uint64, float64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, 0.0, err
	}
	return vmStat.Total, vmStat.Used, vmStat.UsedPercent, nil
}

func getCPUUtilization() (float64, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) == 0 {
		return 0, fmt.Errorf("no cpu percentages returned")
	}
	return percentages[0], nil
}

func getIPAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, intf := range interfaces {
		addrs, err := intf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("no suitable IP address found")
}

func getHostname() (string, error) {
	return os.Hostname()
}

func getCurrentTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func main() {
	fmt.Println("Hello")
	currentTime := getCurrentTime()
	fmt.Printf("Current time: %s\n", currentTime)

	hostname, err := getHostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
	} else {
		fmt.Printf("Hostname: %s\n", hostname)
	}

	ipAddress, err := getIPAddress()
	if err != nil {
		fmt.Printf("Error getting IP address: %v\n", err)
	} else {
		fmt.Printf("IP Address: %s\n", ipAddress)
	}

	cpuUtilization, err := getCPUUtilization()
	if err != nil {
		fmt.Printf("Error getting CPU utilization: %v\n", err)
	} else {
		fmt.Printf("CPU Utilization: %.2f%%\n", cpuUtilization)
	}

	totalMem, usedMem, memPercent, err := getMemoryUtilization()
	if err != nil {
		fmt.Printf("Error getting memory utilization: %v\n", err)
	} else {
		// Simple byte formatting for now
		totalMemGB := float64(totalMem) / (1024 * 1024 * 1024)
		usedMemGB := float64(usedMem) / (1024 * 1024 * 1024)
		fmt.Printf("Memory: Used %.2fGB / Total %.2fGB (%.2f%%)\n", usedMemGB, totalMemGB, memPercent)
	}

	totalDisk, usedDisk, diskPercent, err := getDiskUtilization()
	if err != nil {
		fmt.Printf("Error getting disk utilization: %v\n", err)
	} else {
		totalDiskGB := float64(totalDisk) / (1024 * 1024 * 1024)
		usedDiskGB := float64(usedDisk) / (1024 * 1024 * 1024)
		fmt.Printf("Disk (/): Used %.2fGB / Total %.2fGB (%.2f%%)\n", usedDiskGB, totalDiskGB, diskPercent)
	}

	systemTemp, err := getSystemTemperature()
	if err != nil {
		// For this specific metric, we print N/A on error as requested
		fmt.Printf("System Temperature: N/A (%v)\n", err)
	} else {
		fmt.Printf("System Temperature: %.1f°C\n", systemTemp)
	}
}
