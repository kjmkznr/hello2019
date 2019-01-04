package main

import "fmt"

func main() {
	reqMsg := NewInetDiagReqV2(AF_INET, TCPF_ESTABLISHED)
	respMsg, err := NetlinkInetDiag(reqMsg)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, msg := range respMsg {
		srcIP := msg.SrcIP().To4()
		srcPort := msg.SrcPort()
		dstIP := msg.DstIP().To4()
		dstPort := msg.DstPort()
		fmt.Printf("%v:%v\t%v:%v\t%v\n", srcIP, srcPort, dstIP, dstPort, TCPState(msg.State))
	}
}
