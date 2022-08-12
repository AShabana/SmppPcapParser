package main

import (
	"encoding/binary"
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	MINIMUM_VALID_REQ_ID  uint32 = 1
	MAXIMUM_VALID_REQ_ID  uint32 = 15
	MINIMUM_VALID_RESP_ID uint32 = 2147483649
	MAXIMUM_VALID_RESP_ID uint32 = 2147483665

	BIND_RECEIVER      uint32 = 1
	BIND_RECEIVER_RESP uint32 = 2147483649

	BIND_TRANSMITTER      uint32 = 2
	BIND_TRANSMITTER_RESP uint32 = 2147483650

	BIND_TRANSCEIVER      uint32 = 9
	BIND_TRANSCEIVER_RESP uint32 = 2147483657

	SUBMITSM      uint32 = 4
	SUBMITSM_RESP uint32 = 2147483652

	DELIVERSM      uint32 = 5
	DELIVERSM_RESP uint32 = 2147483653

	UNBIND      uint32 = 6
	UNBIND_RESP uint32 = 2147483654

	MINIMAL_PDU_LEN = 16
)

type PDU_HEADER struct {
	CMDLEN, CMDID, CMDSTATUS, CEQ uint32
}

func extractPduHeader(offset int, payload []byte) *PDU_HEADER {
	if len(payload) > offset+16 && isValidCmd(binary.BigEndian.Uint32(payload[offset+4:offset+8])) {
		return &PDU_HEADER{
			binary.BigEndian.Uint32(payload[offset : offset+4]),     // LEN
			binary.BigEndian.Uint32(payload[offset+4 : offset+8]),   // CMDID
			binary.BigEndian.Uint32(payload[offset+8 : offset+12]),  // STATUS
			binary.BigEndian.Uint32(payload[offset+12 : offset+16]), // CEQ
		}
	}
	if len(payload) > offset+12 && isValidCmd(binary.BigEndian.Uint32(payload[offset+4:offset+8])) {
		return &PDU_HEADER{
			binary.BigEndian.Uint32(payload[offset : offset+4]),    // LEN
			binary.BigEndian.Uint32(payload[offset+4 : offset+8]),  // CMDID
			binary.BigEndian.Uint32(payload[offset+8 : offset+12]), // STATUS
			uint32(0),
		}
	}
	// if len(payload) > 1000 {
	// 	try_traverse(payload)
	// }
	fmt.Printf("! %v", len(payload) > offset+16)
	//time.Sleep(time.Second)
	return nil
}

func isValidCmd(cmdid uint32) bool {
	return MINIMUM_VALID_REQ_ID <= cmdid && cmdid <= MAXIMUM_VALID_REQ_ID || (MINIMUM_VALID_RESP_ID <= cmdid && cmdid <= MAXIMUM_VALID_RESP_ID)
}

func isValidPacket(packet gopacket.Packet) bool {
	return packet != nil && packet.NetworkLayer() != nil && packet.TransportLayer() != nil && packet.TransportLayer().LayerType() == layers.LayerTypeTCP
}

func try_traverse(payload []byte) {
	fmt.Println("XXXXXXXXXXXXX")
	for _, i := range payload {
		fmt.Printf(" %v ", uint8(payload[i]))
		/*if r == int(4) {*/
		//fmt.Println("We expect the current pdu len IN : %v ", payload[i-8:i+57])
		//fmt.Println("We expect the current pdu len is : %v ", binary.BigEndian.Uint32(payload[i-9:i-5]))
		//break
		//}
	}
	fmt.Println()
	fmt.Println("XXXXXXXXXXXXX")
}
