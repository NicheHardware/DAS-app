package main

import (
	"fmt"
	"time"
)

var NodeTypeMap = map[uint16]NodeConst{
	0xDA01: {Type: "Temperature", Unit: "°C", Digit: 1, Color: "#ff8133"},
	0xDA02: {Type: "Voltage", Unit: "V", Digit: 2, Color: "#3377ff"},
	0xDA03: {Type: "Current", Unit: "A", Digit: 3, Color: "#d1dc00"},
	0xDA04: {Type: "Power", Unit: "W", Digit: 1, Color: "#cd0000"},
	0xDA05: {Type: "Frequency", Unit: "Hz", Digit: 1, Color: "#00ff22"},
}

type ScanData struct {
	Mac   [6]uint8
	Leds  uint8
	Type  uint16
	Index uint32

	Data []uint8
}
type NodeConst struct {
	Type  string
	Unit  string
	Digit int
	Color string
}
type Node struct {
	Name string
	NodeConst
	Battery int
	RSSI    int
	Mac     [6]uint8
	Leds    [5]bool
	AxisX   []int64
	AxisY   []float32

	CurrentValue float32
	lastIdx      uint32
	startIdx     uint32
	startTime    time.Time
}

type Nodes map[string]*Node

var nodes Nodes

func init() {
	nodes = make(Nodes)
}

func (ns *Nodes) Clear() {
	*ns = make(Nodes)
}

func (n *Node) InitFrom(data ScanData) {
	n.startTime = time.Now()
	n.startIdx = data.Index
	if _, ok := NodeTypeMap[data.Type]; ok {
		n.Type = NodeTypeMap[data.Type].Type
		n.Unit = NodeTypeMap[data.Type].Unit
		n.Digit = NodeTypeMap[data.Type].Digit
		n.Color = NodeTypeMap[data.Type].Color
	}
	n.Mac = data.Mac
}
func (n *Node) UpdateFrom(data ScanData) {
	// TODO
}

func OnScanData(d ScanData) {
	var macStr string = MacStr(d.Mac)
	if _, ok := nodes[macStr]; !ok {
		nodes[macStr] = &Node{}
		nodes[macStr].InitFrom(d)
	}
	nodes[macStr].UpdateFrom(d)
}
func MacStr(mac [6]byte) string {
	var macStr string
	for i := 0; i < 6; i++ {
		macStr += fmt.Sprintf("%02x", mac[i])
	}
	return macStr
}
