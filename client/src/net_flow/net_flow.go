package net_flow

import (
	"database/sql"
	//"fmt"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
	"util"
	"path"
)

type NetFlowHandler struct {
	db *sql.DB
}

func (h *NetFlowHandler) tryHandleIt(m *nsq.Message) (err error) {
	bodyParts := strings.Split(string(m.Body), "\r\n")

	time_index, err := strconv.Atoi(bodyParts[1])
	netFlowDataParts := strings.Split(bodyParts[5], ",")[1:]
	networkCardNum := len(netFlowDataParts)/4

	beginIndex := 0
	endIndex := networkCardNum
	outBytes := 0
	for index := beginIndex; index < endIndex; index++ {
		oneOutByteFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		outBytes += int(oneOutByteFloat)
	}
	beginIndex = endIndex
	endIndex += networkCardNum
	inBytes := 0
	for index := beginIndex; index < endIndex; index++ {
		oneInByteFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		inBytes += int(oneInByteFloat)
	}

	beginIndex = endIndex
	endIndex += networkCardNum
	outPackets := 0
	for index := beginIndex; index < endIndex; index++ {
		oneOutPacketFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		outPackets += int(oneOutPacketFloat)
	}

	beginIndex = endIndex
	endIndex += networkCardNum
	inPackets := 0
	for index := beginIndex; index < endIndex; index++ {
		oneInPacketsFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		inPackets += int(oneInPacketsFloat)
	}

	sql := `
	INSERT INTO net_flow (date, time_index, ip, host_name, hardware_addr, out_bytes, in_bytes, out_packets, in_packets) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], outBytes, inBytes, outPackets, inPackets)
	return err
}

func (h *NetFlowHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能

	按指标叠加所有网卡的数据
	*/
	defer util.HandleException(path.Join(util.LogRoot, "net_flow.log"), string(m.Body))
	err = h.tryHandleIt(m)
	return err
}

func NewNetFlowHandler(dbLink *sql.DB) (netFlowHandler *NetFlowHandler, err error) {
	netFlowHandler = &NetFlowHandler {
		db: dbLink,
	}
	return netFlowHandler, err
}




