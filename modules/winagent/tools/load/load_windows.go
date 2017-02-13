package load

import "github.com/gaobrian/open-falcon-backend/modules/winagent/tools/wmi"


type LoadData struct {
	LoadPercentage uint64
}

func LoadAvg() (*LoadAvgStat, error) {
	ret := LoadAvgStat{}


	var dst  []LoadData

	err :=  wmi.Query("Select LoadPercentage From Win32_Processor",&dst)
	if err != nil {
		return &ret, err
	}

	ret.Load1 = float64(dst[0].LoadPercentage)  / 100

	return &ret, err
}

