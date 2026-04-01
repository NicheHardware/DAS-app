export interface TypeMetric {
	Key: string
	Name: string
	Unit: string
	Digit: number
	Color: string
	Kind: "source" | "derived"
	CurrentValue: number
	AxisY: number[]
}

export interface TypeNode {
	Name: string
	TypeCode: number
	TypeName: string
	Battery: number
	RSSI: number
	Mac: number[]
	Leds: number[]
	AxisX: number[]
	Metrics: TypeMetric[]
}
