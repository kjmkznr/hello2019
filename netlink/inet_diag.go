package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

// Enums / Constants

// SOCK_DIAG_BY_FAMILY is the netlink message type for requestion socket
// diag data by family. This is newer and can be used with inet_diag_req_v2.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/sock_diag.h#L6
const SOCK_DIAG_BY_FAMILY = 20

// AddressFamily is the address family of the socket.
type AddressFamily uint8

// https://github.com/torvalds/linux/blob/v4.20/include/linux/socket.h#L160
const (
	AF_INET  AddressFamily = 2
	AF_INET6               = 10
)

var addressFamilyNames = map[AddressFamily]string{
	AF_INET:  "ipv4",
	AF_INET6: "ipv6",
}

func (af AddressFamily) String() string {
	if fam, found := addressFamilyNames[af]; found {
		return fam
	}
	return fmt.Sprintf("UNKNOWN (%d)", af)
}

// TCPState represents the state of a TCP connection.
type TCPState uint8

// https://github.com/torvalds/linux/blob/v4.20/include/net/tcp_states.h#L16
const (
	TCP_ESTABLISHED TCPState = iota + 1
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
	TCP_NEW_SYN_RECV
)

var tcpStateNames = map[TCPState]string{
	TCP_ESTABLISHED:  "ESTAB",
	TCP_SYN_SENT:     "SYN-SENT",
	TCP_SYN_RECV:     "SYN-RECV",
	TCP_FIN_WAIT1:    "FIN-WAIT-1",
	TCP_FIN_WAIT2:    "FIN-WAIT-2",
	TCP_TIME_WAIT:    "TIME-WAIT",
	TCP_CLOSE:        "UNCONN",
	TCP_CLOSE_WAIT:   "CLOSE-WAIT",
	TCP_LAST_ACK:     "LAST-ACK",
	TCP_LISTEN:       "LISTEN",
	TCP_CLOSING:      "CLOSING",
	TCP_NEW_SYN_RECV: "NEW-SYN-RECV",
}

func (s TCPState) String() string {
	if state, found := tcpStateNames[s]; found {
		return state
	}
	return "UNKNOWN"
}

const (
	// AllTCPStates is a flag to request all sockets in any TCP state.
	AllTCPStates      = ^uint32(0)
	TCPF_ESTABLISHED  = (1 << TCP_ESTABLISHED)
	TCPF_SYN_SENT     = (1 << TCP_SYN_SENT)
	TCPF_SYN_RECV     = (1 << TCP_SYN_RECV)
	TCPF_FIN_WAIT1    = (1 << TCP_FIN_WAIT1)
	TCPF_FIN_WAIT2    = (1 << TCP_FIN_WAIT2)
	TCPF_TIME_WAIT    = (1 << TCP_TIME_WAIT)
	TCPF_CLOSE        = (1 << TCP_CLOSE)
	TCPF_CLOSE_WAIT   = (1 << TCP_CLOSE_WAIT)
	TCPF_LAST_ACK     = (1 << TCP_LAST_ACK)
	TCPF_LISTEN       = (1 << TCP_LISTEN)
	TCPF_CLOSING      = (1 << TCP_CLOSING)
	TCPF_NEW_SYN_RECV = (1 << TCP_NEW_SYN_RECV)
)

var byteOrder = getEndian()

// NetlinkInetDiag ...
func NetlinkInetDiag(request syscall.NetlinkMessage) ([]*InetDiagMsg, error) {
	s, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_INET_DIAG)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(s)

	lsa := &syscall.SockaddrNetlink{Family: syscall.AF_NETLINK}
	if err := syscall.Sendto(s, serialize(request), 0, lsa); err != nil {
		return nil, err
	}

	// Default size used in libnl.
	readBuf := make([]byte, os.Getpagesize())

	var inetDiagMsgs []*InetDiagMsg
done:
	for {
		buf := readBuf
		nr, _, err := syscall.Recvfrom(s, buf, 0)
		if err != nil {
			return nil, err
		}
		if nr < syscall.NLMSG_HDRLEN {
			return nil, syscall.EINVAL
		}

		buf = buf[:nr]

		msgs, err := syscall.ParseNetlinkMessage(buf)
		if err != nil {
			return nil, err
		}

		for _, m := range msgs {
			if m.Header.Type == syscall.NLMSG_DONE {
				break done
			}
			if m.Header.Type == syscall.NLMSG_ERROR {
				return nil, errors.New("received netlink error (data too short to read errno)")
			}

			inetDiagMsg, err := ParseInetDiagMsg(m.Data)
			if err != nil {
				return nil, err
			}
			inetDiagMsgs = append(inetDiagMsgs, inetDiagMsg)
		}
	}
	return inetDiagMsgs, nil
}

func serialize(msg syscall.NetlinkMessage) []byte {
	msg.Header.Len = uint32(syscall.SizeofNlMsghdr + len(msg.Data))
	b := make([]byte, msg.Header.Len)
	byteOrder.PutUint32(b[0:4], msg.Header.Len)
	byteOrder.PutUint16(b[4:6], msg.Header.Type)
	byteOrder.PutUint16(b[6:8], msg.Header.Flags)
	byteOrder.PutUint32(b[8:12], msg.Header.Seq)
	byteOrder.PutUint32(b[12:16], msg.Header.Pid)
	copy(b[16:], msg.Data)
	return b
}

// V2 Request

var sizeofInetDiagReqV2 = int(unsafe.Sizeof(InetDiagReqV2{}))

// InetDiagReqV2 (inet_diag_req_v2) is used to request diagnostic data.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L37
type InetDiagReqV2 struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	ID       InetDiagSockID
}

func (r InetDiagReqV2) toWireFormat() []byte {
	buf := bytes.NewBuffer(make([]byte, sizeofInetDiagReqV2))
	buf.Reset()
	if err := binary.Write(buf, byteOrder, r); err != nil {
		// This never returns an error.
		panic(err)
	}
	return buf.Bytes()
}

// NewInetDiagReqV2 returns a new NetlinkMessage whose payload is an
// InetDiagReqV2. Callers should set their own sequence number in the returned
// message header.
func NewInetDiagReqV2(af AddressFamily, state uint32) syscall.NetlinkMessage {
	hdr := syscall.NlMsghdr{
		Type:  uint16(SOCK_DIAG_BY_FAMILY),
		Flags: uint16(syscall.NLM_F_DUMP | syscall.NLM_F_REQUEST),
		Pid:   uint32(0),
	}
	req := InetDiagReqV2{
		Family:   uint8(af),
		Protocol: syscall.IPPROTO_TCP,
		States:   state,
	}

	return syscall.NetlinkMessage{Header: hdr, Data: req.toWireFormat()}
}
