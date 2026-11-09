package model

import "github.com/shirou/gopsutil/v3/net"

type IOCountersStat struct {
	Name        string `json:"name"`        // interface name
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	Errin       uint64 `json:"errin"`       // total number of errors while receiving
	Errout      uint64 `json:"errout"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout"`     // total number of FIFO buffers errors while sending
	State       string `json:"state"`
	Time        int64  `json:"time"`
}

// IOCountersFrom copies what gopsutil counted for one interface. It used to be a
// reinterpretation through unsafe.Pointer, which read this struct's length out of
// gopsutil's shorter one: the 24 bytes past its end landed in State and Time, so
// State held whatever pointer sat there until the caller overwrote it, and the
// garbage collector could find it first.
func IOCountersFrom(counters net.IOCountersStat) IOCountersStat {
	return IOCountersStat{
		Name:        counters.Name,
		BytesSent:   counters.BytesSent,
		BytesRecv:   counters.BytesRecv,
		PacketsSent: counters.PacketsSent,
		PacketsRecv: counters.PacketsRecv,
		Errin:       counters.Errin,
		Errout:      counters.Errout,
		Dropin:      counters.Dropin,
		Dropout:     counters.Dropout,
		Fifoin:      counters.Fifoin,
		Fifoout:     counters.Fifoout,
	}
}
