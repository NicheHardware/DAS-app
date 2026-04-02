package main

import (
	"fmt"
	"sort"
	"time"
)

type MetricKind string

const (
	MetricSource  MetricKind = "source"
	MetricDerived MetricKind = "derived"
)

type MetricDef struct {
	Key   string
	Name  string
	Unit  string
	Digit int
	Color string
	Kind  MetricKind
}

type TickToMillisFunc func(deltaTick uint64) int64
type ComputeMetricFunc func(values map[string]float32) float32

type NodeTypeDef struct {
	Code         uint16
	Name         string
	Metrics      []MetricDef
	TickToMillis TickToMillisFunc
	Compute      map[string]ComputeMetricFunc
}

var NodeTypeMap = map[uint16]NodeTypeDef{
	0xDA01: {
		Code:         0xDA01,
		Name:         "Temperature",
		Metrics:      []MetricDef{{Key: "temperature", Name: "Temperature", Unit: "°C", Digit: 1, Color: "#ff8133", Kind: MetricSource}},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
	},
	0xDA02: {
		Code:         0xDA02,
		Name:         "Voltage",
		Metrics:      []MetricDef{{Key: "voltage", Name: "Voltage", Unit: "V", Digit: 2, Color: "#3377ff", Kind: MetricSource}},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
	},
	0xDA03: {
		Code:         0xDA03,
		Name:         "Current",
		Metrics:      []MetricDef{{Key: "current", Name: "Current", Unit: "A", Digit: 3, Color: "#d1dc00", Kind: MetricSource}},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
	},
	0xDA04: {
		Code:         0xDA04,
		Name:         "Power",
		Metrics:      []MetricDef{{Key: "power", Name: "Power", Unit: "W", Digit: 1, Color: "#cd0000", Kind: MetricSource}},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
	},
	0xDA05: {
		Code:         0xDA05,
		Name:         "Frequency",
		Metrics:      []MetricDef{{Key: "frequency", Name: "Frequency", Unit: "Hz", Digit: 1, Color: "#00ff22", Kind: MetricSource}},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
	},
	0xDC01: {
		Code: 0xDC01,
		Name: "CM01",
		Metrics: []MetricDef{
			{Key: "voltage", Name: "Voltage", Unit: "V", Digit: 1, Color: "#3377ff", Kind: MetricSource},
			{Key: "current", Name: "Current", Unit: "A", Digit: 1, Color: "#d1dc00", Kind: MetricSource},
			{Key: "power", Name: "Power", Unit: "W", Digit: 1, Color: "#cd0000", Kind: MetricDerived},
		},
		TickToMillis: func(deltaTick uint64) int64 { return int64(deltaTick) },
		Compute: map[string]ComputeMetricFunc{
			"power": func(values map[string]float32) float32 {
				return values["voltage"] * values["current"]
			},
		},
	},
}

type ScanData struct {
	Mac   [6]uint8
	Leds  uint8
	Type  uint16
	Index uint32
	Tick  uint64
	Data  []uint8
}

type MetricSeries struct {
	Key          string
	Name         string
	Unit         string
	Digit        int
	Color        string
	Kind         MetricKind
	CurrentValue float32
	AxisY        []float32
}

type Node struct {
	Name     string
	TypeCode uint16
	TypeName string
	Battery  int
	RSSI     int
	Mac      [6]uint8
	Leds     [5]bool
	AxisX    []int64
	Metrics  []*MetricSeries

	FirstTick    uint64
	FirstLocalTs int64
	TimeSynced   bool
	lastIdx      uint32
	startIdx     uint32
	startTime    time.Time
	LastDeviceTs uint64
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
	n.TypeCode = data.Type
	if def, ok := NodeTypeMap[data.Type]; ok {
		n.TypeName = def.Name
		n.Metrics = make([]*MetricSeries, 0, len(def.Metrics))
		for _, metricDef := range def.Metrics {
			n.Metrics = append(n.Metrics, &MetricSeries{
				Key:   metricDef.Key,
				Name:  metricDef.Name,
				Unit:  metricDef.Unit,
				Digit: metricDef.Digit,
				Color: metricDef.Color,
				Kind:  metricDef.Kind,
			})
		}
	}
	n.Mac = data.Mac
}

func (n *Node) NormalizeTimestamp(tick uint64, hostNow int64) int64 {
	def, ok := NodeTypeMap[n.TypeCode]
	if !ok || def.TickToMillis == nil {
		return hostNow
	}
	if !n.TimeSynced {
		n.FirstTick = tick
		n.FirstLocalTs = hostNow
		n.TimeSynced = true
		return hostNow
	}
	return n.FirstLocalTs + def.TickToMillis(tick-n.FirstTick)
}

func (n *Node) AppendFrame(ts int64, values map[string]float32) {
	if len(n.AxisX) > 0 && n.AxisX[len(n.AxisX)-1] == ts {
		for _, metric := range n.Metrics {
			value := values[metric.Key]
			metric.CurrentValue = value
			if len(metric.AxisY) > 0 {
				metric.AxisY[len(metric.AxisY)-1] = value
			}
		}
		return
	}
	n.AxisX = append(n.AxisX, ts)
	for _, metric := range n.Metrics {
		value := values[metric.Key]
		metric.CurrentValue = value
		metric.AxisY = append(metric.AxisY, value)
	}
	if len(n.AxisX) > 512 {
		n.AxisX = n.AxisX[len(n.AxisX)-512:]
		for _, metric := range n.Metrics {
			if len(metric.AxisY) > 512 {
				metric.AxisY = metric.AxisY[len(metric.AxisY)-512:]
			}
		}
	}
}

func (n *Node) UpdateFrom(data ScanData) {
	n.lastIdx = data.Index
	_ = data
}

func (n *Node) MetricMap() map[string]*MetricSeries {
	metrics := make(map[string]*MetricSeries, len(n.Metrics))
	for _, metric := range n.Metrics {
		metrics[metric.Key] = metric
	}
	return metrics
}

func (n *Node) ApplyMetricRecords(records map[uint64]map[string]float32, hostNow int64) bool {
	if len(records) == 0 {
		return false
	}
	keys := make([]uint64, 0, len(records))
	for ts := range records {
		if ts <= n.LastDeviceTs {
			continue
		}
		keys = append(keys, ts)
	}
	if len(keys) == 0 {
		return false
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, tick := range keys {
		values := records[tick]
		if def, ok := NodeTypeMap[n.TypeCode]; ok {
			for key, compute := range def.Compute {
				values[key] = compute(values)
			}
		}
		ts := n.NormalizeTimestamp(tick, hostNow)
		n.AppendFrame(ts, values)
		n.LastDeviceTs = tick
	}
	return true
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
