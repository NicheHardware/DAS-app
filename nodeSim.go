package main

import (
	"fmt"
	"math/rand"
	"time"
)

type SimState struct {
	Type      uint16
	Delay     int
	Current   float32
	Target    float32
	TargetMax float32
	TargetMin float32
	Delta     float32
	DeltaMax  float32
	DeltaMin  float32
	// speed = delta * k
	k float32
}
type SimStates map[string]*SimState

var simStates SimStates

func init() {
	simStates = make(SimStates)
}

func (ss *SimStates) Clear() {
	*ss = make(SimStates)
}

func (ss *SimState) init(tp uint16) {
	switch tp {
	case 0xDA01:
		ss.k = 1
		ss.TargetMax = 110
		ss.TargetMin = 18
	case 0xDA02:
		ss.k = 5
		switch rand.Intn(3) {
		case 0:
			ss.TargetMax = 13
			ss.TargetMin = 11
		case 1:
			ss.TargetMax = 5.5
			ss.TargetMin = 4.5
		case 2:
			ss.TargetMax = 3.6
			ss.TargetMin = 2.7
		}
	case 0xDA03:
		ss.k = 5
		ss.TargetMax = 33
		ss.TargetMin = 2
	case 0xDA04:
		ss.k = 5
		ss.TargetMax = 600
		ss.TargetMin = 50
	case 0xDA05:
		ss.k = 2
		switch rand.Intn(2) {
		case 0:
			ss.TargetMax = 52
			ss.TargetMin = 48
		case 1:
			ss.TargetMax = 8000
			ss.TargetMin = 800
		}
	}

	ss.DeltaMax = (ss.TargetMax - ss.TargetMin) * 0.02 * ss.k
	ss.DeltaMin = -(ss.TargetMax - ss.TargetMin) * 0.02 * ss.k
	ss.Target = (ss.TargetMax-ss.TargetMin)*rand.Float32() + ss.TargetMin
	ss.Current = ss.Target
}

func (ss *SimState) update() {
	valueRange := ss.TargetMax - ss.TargetMin
	deltaRange := ss.DeltaMax - ss.DeltaMin

	if ss.Target > ss.Current+deltaRange {
		ss.Delta += deltaRange * (rand.Float32() - 0.4) / 10
	} else if ss.Target < ss.Current-deltaRange {
		ss.Delta += deltaRange * (rand.Float32() - 0.6) / 10
	} else {
		ss.Delay = rand.Intn(15)
		ss.Target = valueRange*rand.Float32() + ss.TargetMin
	}

	if ss.Delta > ss.DeltaMax+deltaRange {
		ss.Delta = ss.DeltaMax + deltaRange
	}
	if ss.Delta < ss.DeltaMin-deltaRange {
		ss.Delta = ss.DeltaMin - deltaRange
	}

	if ss.Delay > 0 {
		ss.Delay--
		ss.Delta = 0
		ss.Current += deltaRange * (rand.Float32() - 0.5) / 10
	} else {
		ss.Current += ss.Delta
	}

	if ss.Current > ss.TargetMax+valueRange/10 {
		ss.Current = ss.TargetMax + valueRange/10
	}
	if ss.Current < ss.TargetMin-valueRange/10 {
		ss.Current = ss.TargetMin - valueRange/10
	}
}

func newSimNode(tp uint16) (*Node, *SimState) {
	node := &Node{}
	node.InitFrom(ScanData{
		Index: uint32(rand.Intn(65536)),
		Type:  tp,
		Mac: [6]byte{
			byte(rand.Intn(0x100)),
			byte(rand.Intn(0x100)),
			byte(rand.Intn(0x100)),
			byte(rand.Intn(0x100)),
			byte(rand.Intn(0x100)),
			byte(rand.Intn(0x100))},
	})
	node.Name = fmt.Sprintf("Node %03d%03d", rand.Intn(256), rand.Intn(1000))
	node.Battery = 15 + rand.Intn(80)
	node.RSSI = -30 - rand.Intn(60)
	rled := rand.Intn(16) + 1
	for i := range 5 {
		node.Leds[i] = rled&(1<<i) > 0
	}

	simState := &SimState{}
	simState.init(tp)
	return node, simState
}

func generateSimNodes(count int) {
	for range count {
		node, simState := newSimNode(uint16(rand.Intn(5) + 0xDA01))
		macStr := MacStr(node.Mac)
		simStates[macStr] = simState
		nodes[macStr] = node
	}
}

var simulating bool = false

func StartSim(count int) {
	generateSimNodes(count)
	simulating = true
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			if !simulating {
				break
			}
			updateSimNodes()
			app.NewDataNotify()
		}
	}()
}
func StopSim() {
	simulating = false
	nodes.Clear()
	simStates.Clear()
}
func updateSimNodes() {
	for _, node := range nodes {
		if simState, ok := simStates[MacStr(node.Mac)]; ok {
			simState.update()
			node.CurrentValue = float32(simState.Current)
			node.AxisY = append(node.AxisY, float32(simState.Current))
			node.AxisX = append(node.AxisX, time.Now().UnixMilli())
		}
	}
}
