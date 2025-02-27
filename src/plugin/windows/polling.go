package windows

import (
	"NMS/src/util"
	"bytes"
	"errors"
	"fmt"
	"github.com/masterzen/winrm"
	"os"

	"strings"
	"time"
)

type Config struct {
	IP       string
	Username string
	Password string
	timeout  time.Duration
}

var (
	client *winrm.Client

	shell *winrm.Shell

	logInstance = util.InitializeLogger()
)

func initWinRMClient(config Config) error {

	endpoint := winrm.NewEndpoint(config.IP, 5985, false, false, nil, nil, nil, config.timeout)

	var err error

	client, err = winrm.NewClient(endpoint, config.Username, config.Password)

	if err != nil {

		logInstance.LogError(fmt.Errorf("failed to create WinRM client: %v", err))

		return err

	}

	return nil
}

func initWinRMShell() error {

	var err error

	shell, err = client.CreateShell()

	if err != nil {

		logInstance.LogError(fmt.Errorf("Failed to create WinRM shell: %v", err))

		return err

	}

	return nil

}

func closeWinRMShell() {

	if shell != nil {

		shell.Close()

		shell = nil

		logInstance.LogInfo("WinRM shell closed")

	}

}

func executeAndFetchWindowsCounters(command string) string {

	if client == nil {

		logInstance.LogError(errors.New("WinRM client is not initialized"))

		return "0"

	}

	var stdout, stderr bytes.Buffer

	exitCode, err := client.Run("powershell -ExecutionPolicy Bypass -Command \""+command+"\"", &stdout, &stderr)

	if err != nil || exitCode != 0 {

		logInstance.LogError(fmt.Errorf("Command failed: %s, Error: %v, Stderr: %s", command, err, stderr.String()))

		return "0"
	}

	output := strings.TrimSpace(stdout.String())

	return output

}

func putDataIntoCollection(res map[string]interface{}, outputLines []string) (map[string]interface{}, error) {

	res[util.SystemHostName] = strings.TrimSpace(outputLines[0])

	res[util.SystemUpTime] = strings.TrimSpace(outputLines[1])

	res[util.SystemDiskUsedBytes] = strings.TrimSpace(outputLines[2])

	res[util.SystemPhysicalProcessors] = strings.TrimSpace(outputLines[3])

	res[util.SystemCPUCores] = strings.TrimSpace(outputLines[4])

	res[util.SystemLogicalProcessors] = strings.TrimSpace(outputLines[5])

	res[util.SystemRunningProcesses] = strings.TrimSpace(outputLines[6])

	res[util.SystemOSVersion] = strings.TrimSpace(outputLines[7])

	res[util.SystemVendor] = strings.TrimSpace(outputLines[8])

	res[util.SystemSerialNumber] = strings.TrimSpace(outputLines[9])

	res[util.SystemCPUIdlePercent] = strings.TrimSpace(outputLines[10])

	res[util.SystemMemoryFreePercent] = strings.TrimSpace(outputLines[11])

	res[util.SystemCacheMemoryBytes] = strings.TrimSpace(outputLines[12])

	res[util.SystemMemoryUsedPercent] = strings.TrimSpace(outputLines[13])

	res[util.SystemMemoryAvailableBytes] = strings.TrimSpace(outputLines[14])

	res[util.SystemCPUDescription] = strings.TrimSpace(outputLines[15])

	res[util.SystemCPUInterruptPerSec] = strings.TrimSpace(outputLines[16])

	res[util.SystemMemoryCommittedBytes] = strings.TrimSpace(outputLines[17])

	res[util.SystemDiskFreePercent] = strings.TrimSpace(outputLines[18])

	res[util.SystemDiskUsedPercent] = strings.TrimSpace(outputLines[19])

	res[util.SystemNetworkTCPConnections] = strings.TrimSpace(outputLines[20])

	res[util.SystemContextSwitchesPerSec] = strings.TrimSpace(outputLines[21])

	res[util.SystemDiskCapacityBytes] = strings.TrimSpace(outputLines[22])

	res[util.SystemCPUType] = strings.TrimSpace(outputLines[23])

	res[util.SystemName] = strings.TrimSpace(outputLines[24])

	res[util.SystemThreads] = strings.TrimSpace(outputLines[25])

	res[util.SystemProcessorQueueLength] = strings.TrimSpace(outputLines[26])

	res[util.SystemCPUUserPercent] = strings.TrimSpace(outputLines[27])

	res[util.SystemCPUPercent] = strings.TrimSpace(outputLines[28])

	res[util.SystemMemoryInstalledBytes] = strings.TrimSpace(outputLines[29])

	res[util.SystemMemoryUsedBytes] = strings.TrimSpace(outputLines[30])

	res[util.SystemDiskFreeBytes] = strings.TrimSpace(outputLines[31])

	res[util.SystemMemoryFreeBytes] = strings.TrimSpace(outputLines[32])

	return res, nil

}

func Start(ip string, username string, password string) map[string]interface{} {

	res := make(map[string]interface{})

	config := Config{

		IP: ip,

		Username: username,

		Password: password,

		timeout: 2 * time.Minute,
	}

	if err := initWinRMClient(config); err != nil {

		logInstance.LogInfo(fmt.Sprintf("Failed to initialize WinRM client: %v", err))

	}

	if err := initWinRMShell(); err != nil {

		logInstance.LogError(fmt.Errorf("Error creating WinRM shell: %v", err))

		os.Exit(1)

	}

	defer closeWinRMShell()

	psScript := `$os = Get-CimInstance Win32_OperatingSystem;$processor = Get-CimInstance Win32_Processor; $memory = Get-CimInstance Win32_PerfFormattedData_PerfOS_Memory; $disk = Get-CimInstance Win32_LogicalDisk; $bios = Get-CimInstance Win32_BIOS; $cpuPerf = Get-Counter '\Processor(_Total)\% Idle Time' | Select-Object -ExpandProperty CounterSamples | Select-Object -ExpandProperty CookedValue;$env:COMPUTERNAME;($os.LastBootUpTime - (Get-Date)).TotalSeconds; ($disk | Measure-Object -Property Size -Sum).Sum - ($disk | Measure-Object -Property FreeSpace -Sum).Sum;(Get-CimInstance Win32_ComputerSystem).NumberOfProcessors;($processor | Measure-Object -Property NumberOfCores -Sum).Sum; (Get-CimInstance Win32_ComputerSystem).NumberOfLogicalProcessors;(Get-Process | Measure-Object).Count;$os.Caption;(Get-CimInstance Win32_ComputerSystem).Manufacturer;$bios.SerialNumber;[math]::Round($cpuPerf, 2);[math]::Round(($os.FreePhysicalMemory / $os.TotalVisibleMemorySize) * 100, 2); $memory.CacheBytes; [math]::Round((($os.TotalVisibleMemorySize - $os.FreePhysicalMemory) / $os.TotalVisibleMemorySize) * 100, 2); $os.FreePhysicalMemory * 1024; $processor.Name; (Get-Counter '\Processor(_Total)\Interrupts/sec').CounterSamples.CookedValue; ($os.TotalVirtualMemorySize - $os.FreeVirtualMemory) * 1024;[math]::Round(($disk | Measure-Object -Property FreeSpace -Sum).Sum * 100 / ($disk | Measure-Object -Property Size -Sum).Sum, 2);[math]::Round((($disk | Measure-Object -Property Size -Sum).Sum - ($disk | Measure-Object -Property FreeSpace -Sum).Sum) * 100 / ($disk | Measure-Object -Property Size -Sum).Sum, 2);(Get-CimInstance Win32_PerfRawData_Tcpip_TCPv4).ConnectionsEstablished;(Get-CimInstance Win32_PerfFormattedData_PerfOS_System).ContextSwitchesPerSec;($disk | Measure-Object -Property Size -Sum).Sum; $processor.Name; $env:COMPUTERNAME;(Get-Process | ForEach-Object { $_.Threads.Count } | Measure-Object -Sum).Sum; (Get-CimInstance Win32_PerfRawData_PerfOS_System).ProcessorQueueLength;(Get-Counter '\Processor(_Total)\% User Time').CounterSamples.CookedValue;(Get-Counter '\Processor(_Total)\% Processor Time').CounterSamples.CookedValue;$os.TotalVisibleMemorySize * 1024;(($os.TotalVisibleMemorySize - $os.FreePhysicalMemory) * 1024);($disk | Measure-Object -Property FreeSpace -Sum).Sum;$os.FreePhysicalMemory * 1024;`

	output := executeAndFetchWindowsCounters(psScript)

	outputLines := strings.Split(output, "\r\n")

	res, _ = putDataIntoCollection(res, outputLines)

	return res

}
