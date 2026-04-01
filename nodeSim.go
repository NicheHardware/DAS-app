package main

import (
	"fmt"
	"math/rand"
	"time"
)

type SimMetricState struct {
	Delay     int
	Current   float32
	Target    float32
	TargetMax float32
	TargetMin float32
	Delta     float32
	DeltaMax  float32
	DeltaMin  float32
	k         float32
}

type SimState struct {
	Type    uint16
	Tick    uint64
	Metrics map[string]*SimMetricState
}

type SimStates map[string]*SimState

var simStates SimStates

func init() {
	simStates = make(SimStates)
}

func (ss *SimStates) Clear() {
	*ss = make(SimStates)
}

func newMetricState(key string) *SimMetricState {
	ms := &SimMetricState{}
	switch key {
	case "temperature":
		ms.k = 1
		ms.TargetMax = 110
		ms.TargetMin = 18
	case "voltage":
		ms.k = 5
		switch rand.Intn(3) {
		case 0:
			ms.TargetMax = 13
			ms.TargetMin = 11
		case 1:
			ms.TargetMax = 5.5
			ms.TargetMin = 4.5
		default:
			ms.TargetMax = 3.6
			ms.TargetMin = 2.7
		}
	case "current":
		ms.k = 5
		ms.TargetMax = 33
		ms.TargetMin = 2
	case "power":
		ms.k = 5
		ms.TargetMax = 600
		ms.TargetMin = 50
	case "frequency":
		ms.k = 2
		switch rand.Intn(2) {
		case 0:
			ms.TargetMax = 52
			ms.TargetMin = 48
		default:
			ms.TargetMax = 8000
			ms.TargetMin = 800
		}
	default:
		ms.k = 1
		ms.TargetMax = 100
		ms.TargetMin = 0
	}

	ms.DeltaMax = (ms.TargetMax - ms.TargetMin) * 0.02 * ms.k
	ms.DeltaMin = -(ms.TargetMax - ms.TargetMin) * 0.02 * ms.k
	ms.Target = (ms.TargetMax-ms.TargetMin)*rand.Float32() + ms.TargetMin
	ms.Current = ms.Target
	return ms
}

func (ms *SimMetricState) update() {
	valueRange := ms.TargetMax - ms.TargetMin
	deltaRange := ms.DeltaMax - ms.DeltaMin

	if ms.Target > ms.Current+deltaRange {
		ms.Delta += deltaRange * (rand.Float32() - 0.4) / 10
	} else if ms.Target < ms.Current-deltaRange {
		ms.Delta += deltaRange * (rand.Float32() - 0.6) / 10
	} else {
		ms.Delay = rand.Intn(15)
		ms.Target = valueRange*rand.Float32() + ms.TargetMin
	}

	if ms.Delta > ms.DeltaMax+deltaRange {
		ms.Delta = ms.DeltaMax + deltaRange
	}
	if ms.Delta < ms.DeltaMin-deltaRange {
		ms.Delta = ms.DeltaMin - deltaRange
	}

	if ms.Delay > 0 {
		ms.Delay--
		ms.Delta = 0
		ms.Current += deltaRange * (rand.Float32() - 0.5) / 10
	} else {
		ms.Current += ms.Delta
	}

	if ms.Current > ms.TargetMax+valueRange/10 {
		ms.Current = ms.TargetMax + valueRange/10
	}
	if ms.Current < ms.TargetMin-valueRange/10 {
		ms.Current = ms.TargetMin - valueRange/10
	}
}

func newSimNode(tp uint16) (*Node, *SimState) {
	def := NodeTypeMap[tp]
	node := &Node{}
	node.InitFrom(ScanData{
		Index: uint32(rand.Intn(65536)),
		Type:  tp,
		Tick:  0,
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

	simState := &SimState{Type: tp, Metrics: make(map[string]*SimMetricState, len(def.Metrics))}
	for _, metric := range def.Metrics {
		if metric.Kind == MetricSource {
			simState.Metrics[metric.Key] = newMetricState(metric.Key)
		}
	}
	return node, simState
}

func generateSimNodes(count int) {
	types := []uint16{0xDA01, 0xDA02, 0xDA03, 0xDA04, 0xDA05, 0xDC01}
	for range count {
		node, simState := newSimNode(types[rand.Intn(len(types))])
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
			simState.Tick += 500
			values := map[string]float32{}
			for key, metricState := range simState.Metrics {
				metricState.update()
				values[key] = metricState.Current
			}
			if def, ok := NodeTypeMap[node.TypeCode]; ok {
				for key, compute := range def.Compute {
					values[key] = compute(values)
				}
			}
			ts := node.NormalizeTimestamp(simState.Tick, time.Now().UnixMilli())
			node.AppendFrame(ts, values)
		}
	}
}
