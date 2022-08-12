package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	pcapFile string = "./th500_ws500.pcap"
	//pcapFile string = "./notabletoparse.pcap"
	//pcapFile string = "./sampleSubmits.pcap"
	//pcapFile string = "./sampleResponses.pcap"
	//pcapFile string = "/Users/shabana/Downloads/submit_response_sample.pcap"
	handle *pcap.Handle
	err    error
)

/////// NEW FILE

/////// NEW FILE

func main() {
	// Open file instead of device
	handle, err = pcap.OpenOffline(pcapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter string = "port 10003"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Filter set to .", filter)
	fmt.Println("Pcap file: ", pcapFile)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	var stats = make(map[uint32]map[string]int, len(packets))
	var aggr = make(map[uint32]int, len(packets))

	for packet := range packets {
		fmt.Printf(".")
		packet_time := fmt.Sprintf("%v:%v:%v", packet.Metadata().Timestamp.Hour(), packet.Metadata().Timestamp.Minute(), packet.Metadata().Timestamp.Second())

		if packet == nil {
			return
		}

		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		tcp := packet.TransportLayer().(*layers.TCP)
		payload_length := len(tcp.Payload)
		pduh := extractPduHeader(0, tcp.Payload)
		if pduh == nil {
			continue
		}
		offset := int(pduh.CMDLEN)
		longpacketloopstarted := false
		for payload_length >= offset && isValidCmd(pduh.CMDID) {
			fmt.Printf("*")
			if !longpacketloopstarted {
				offset = 0
			}
			longpacketloopstarted = true
			pduh = extractPduHeader(offset, tcp.Payload)
			if pduh == nil {
				//fmt.Println("Failed to parse packet")
				break
			}
			fmt.Printf("%+v\n", pduh)

			offset = offset + int(pduh.CMDLEN)
			aggr[pduh.CMDID]++
			if _, ok := stats[pduh.CMDID]; !ok {
				stats[pduh.CMDID] = make(map[string]int)
			}
			if _, ok := stats[pduh.CMDID][packet_time]; !ok {
				stats[pduh.CMDID][packet_time] = 0
			}
			stats[pduh.CMDID][packet_time]++
		}
		// fmt.Println("[DEBUG] END Multi PDU Packet traversal")

		// stats[pduh.CMDID][packet_time]++
		// aggr[pduh.CMDID]++

		// switch pduh.CMDID {
		// case BIND_RECEIVER:
		// 	fmt.Println("Bind Receiver")
		// case BIND_RECEIVER_RESP:
		// 	fmt.Println("Bind Receiver Resp")
		// case BIND_TRANSMITTER:
		// 	fmt.Println("Bind Transmitter")
		// case BIND_TRANSMITTER_RESP:
		// 	fmt.Println("Bind Transmitter Resp")
		// case BIND_TRANSCEIVER:
		// 	fmt.Println("Bind Transceiver ")
		// case BIND_TRANSCEIVER_RESP:
		// 	fmt.Println("Bind transceiver Resp")
		// case SUBMITSM:
		// 	fmt.Println("SubmitSM")
		// 	fmt.Println(aggr[pduh.CMDID])
		// case SUBMITSM_RESP:
		// 	fmt.Println("Submit SM Resp")
		// 	fmt.Println(aggr[pduh.CMDID])
		// case uint32(DELIVERSM):
		// 	fmt.Println("Devliver SM")
		// case uint32(DELIVERSM_RESP):
		// 	fmt.Println("DeliverSM Resp")
		// case uint32(UNBIND):
		// 	fmt.Println("Unbind")
		// case uint32(UNBIND_RESP):
		// 	fmt.Println("Unbind resp")
		// default:
		// 	fmt.Println("UNKOWN ", pduh.CMDID)
		// }

	}

	fmt.Println("--------------------------------")
	fmt.Println(aggr)
	fmt.Println("--------------------------------")
	fmt.Println(stats)

	fmt.Println("Total Submit SM")
	for s, v := range stats[SUBMITSM] {
		fmt.Println("second ", s)
		fmt.Println("total sent", v)

	}
	fmt.Println("Total Submit SM Responses")
	for s, v := range stats[SUBMITSM_RESP] {
		fmt.Println("second ", s)
		fmt.Println("total sent", v)

	}
}
