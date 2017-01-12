package net

import (
	"net"
	"os"
	"syscall"
	"unsafe"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/wmi"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/internal/common"
)

var (
	modiphlpapi             = syscall.NewLazyDLL("iphlpapi.dll")
	procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")
	procGetExtendedUdpTable = modiphlpapi.NewProc("GetExtendedUdpTable")
)

const (
	TCPTableBasicListener = iota
	TCPTableBasicConnections
	TCPTableBasicAll
	TCPTableOwnerPIDListener
	TCPTableOwnerPIDConnections
	TCPTableOwnerPIDAll
	TCPTableOwnerModuleListener
	TCPTableOwnerModuleConnections
	TCPTableOwnerModuleAll
)

func NetIOCounters(pernic bool) ([]NetIOCountersStat, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ai, err := getAdapterList()
	if err != nil {
		return nil, err
	}
	var ret []NetIOCountersStat

	for _, ifi := range ifs {
		name := ifi.Name
		for ; ai != nil; ai = ai.Next {
			c := NetIOCountersStat{
				Name: name,
			}

			row := syscall.MibIfRow{Index: ai.Index}
			e := syscall.GetIfEntry(&row)
			if e != nil {
				return nil, os.NewSyscallError("GetIfEntry", e)
			}
			c.BytesSent = uint64(row.OutOctets)
			c.BytesRecv = uint64(row.InOctets)
			c.PacketsSent = uint64(row.OutUcastPkts)
			c.PacketsRecv = uint64(row.InUcastPkts)
			c.Errin = uint64(row.InErrors)
			c.Errout = uint64(row.OutErrors)
			c.Dropin = uint64(row.InDiscards)
			c.Dropout = uint64(row.OutDiscards)

			ret = append(ret, c)
		}
	}

	if pernic == false {
		return getNetIOCountersAll(ret)
	}
	return ret, nil
}

// Return a list of network connections opened by a process
func NetConnections(kind string) ([]NetConnectionStat, error) {
	var ret []NetConnectionStat

	return ret, common.NotImplementedError
}

// borrowed from src/pkg/net/interface_windows.go
func getAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

type Tcpipdatastat struct {
	ConFailures    uint64 `json:"confailures"`
	ConActive      uint64 `json:"conactive"`
	ConPassive     uint64 `json:"conpassive"`
	ConEstablished uint64 `json:"conestablished"`
	ConReset       uint64 `json:"conreset"`
}

type Win32_TCPPerfFormattedData struct {
	ConnectionFailures     uint64
	ConnectionsActive      uint64
	ConnectionsPassive     uint64
	ConnectionsEstablished uint64
	ConnectionsReset       uint64
}

func TcpipCounters() ([]Tcpipdatastat, error) {
	ret := make([]Tcpipdatastat, 0)
	var dst []Win32_TCPPerfFormattedData
	err := wmi.Query("SELECT ConnectionFailures,ConnectionsActive,ConnectionsPassive,ConnectionsEstablished,ConnectionsReset FROM Win32_PerfRawData_Tcpip_TCPv4", &dst)
	if err != nil {
		return ret, err
	}

	for _, d := range dst {

		ret = append(ret, Tcpipdatastat{
			ConFailures:    uint64(d.ConnectionFailures),
			ConActive:      uint64(d.ConnectionsActive),
			ConPassive:     uint64(d.ConnectionsPassive),
			ConEstablished: uint64(d.ConnectionsEstablished),
			ConReset:       uint64(d.ConnectionsReset),
		})
	}

	return ret, nil
}

