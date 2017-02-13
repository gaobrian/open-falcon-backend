package os

import "github.com/gaobrian/open-falcon-backend/modules/winagent/tools/wmi"

type Uptime struct {
	SystemUpTime uint64
}

type UptimeResult struct {
	Days uint64
	Hours uint64
	Mins  uint64
}


func SysupTime2Int64()  (uint64){

	var dst []Uptime
	err := wmi.Query("Select SystemUpTime From Win32_PerfFormattedData_PerfOS_System", &dst)
	if err != nil {
		return 0
	}
	return dst[0].SystemUpTime
}

func QueryUptime()  (t UptimeResult, err error){

	var dst []Uptime
	err = wmi.Query("Select SystemUpTime From Win32_PerfFormattedData_PerfOS_System", &dst)
	if err != nil {
		return
	}

	sysup := Uptime{}
	for _, d := range dst {
		sysup = d
	}

	minTotal := sysup.SystemUpTime / 60.0
	hourTotal := minTotal / 60.0

	days := uint64(hourTotal / 24.0)
	hours := uint64(hourTotal) - days*24
	mins := uint64(minTotal) - (days * 60 * 24) - (hours * 60)

	t.Days = days
	t.Hours = hours
	t.Mins = mins

	return
}

type WinVersion struct {
	CSDVersion string
	Caption string
}

func OSVersion() (ver string, err error) {
	var dst []WinVersion

	err =  wmi.Query("Select CSDVersion,Caption From Win32_OperatingSystem",&dst)
	if err != nil {
		return
	}

	ver = dst[0].Caption +"," + dst[0].CSDVersion
	return
}