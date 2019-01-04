package main

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/pkg/errors"
)

// Response messages.

// InetDiagMsg (inet_diag_msg) is the base info structure. It contains socket
// identity (addrs/ports/cookie) and the information shown by netstat.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L86
type InetDiagMsg struct {
	Family  uint8 // Address family.
	State   uint8 // TCP State
	Timer   uint8
	Retrans uint8

	ID InetDiagSockID

	Expires uint32
	RQueue  uint32 // Recv-Q
	WQueue  uint32 // Send-Q
	UID     uint32 // UID
	Inode   uint32 // Inode of socket.
}

// ParseInetDiagMsg parse an InetDiagMsg from a byte slice. It assumes the
// InetDiagMsg starts at the beginning of b. Invoke this method to parse the
// payload of a netlink response.
func ParseInetDiagMsg(b []byte) (*InetDiagMsg, error) {
	r := bytes.NewReader(b)
	inetDiagMsg := &InetDiagMsg{}
	err := binary.Read(r, byteOrder, inetDiagMsg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal inet_diag_msg")
	}
	return inetDiagMsg, nil
}

// SrcPort returns the source (local) port.
func (m InetDiagMsg) SrcPort() int { return int(binary.BigEndian.Uint16(m.ID.SPort[:])) }

// DstPort returns the destination (remote) port.
func (m InetDiagMsg) DstPort() int { return int(binary.BigEndian.Uint16(m.ID.DPort[:])) }

// SrcIP returns the source (local) IP.
func (m InetDiagMsg) SrcIP() net.IP { return ip(m.ID.Src, AddressFamily(m.Family)) }

// DstIP returns the destination (remote) IP.
func (m InetDiagMsg) DstIP() net.IP { return ip(m.ID.Dst, AddressFamily(m.Family)) }

func ip(data [16]byte, af AddressFamily) net.IP {
	if af == AF_INET {
		return net.IPv4(data[0], data[1], data[2], data[3])
	}
	return net.IP(data[:])
}

// InetDiagSockID (inet_diag_sockid) contains the socket identity.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L13
type InetDiagSockID struct {
	SPort  [2]byte  // Source port (big-endian).
	DPort  [2]byte  // Destination port (big-endian).
	Src    [16]byte // Source IP
	Dst    [16]byte // Destination IP
	If     uint32
	Cookie [2]uint32
}
