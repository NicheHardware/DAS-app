export interface TypeNode {
	Name: string

	Type: string
	Unit: string
	Digit: number
	Color: string

	Battery: number
	RSSI: number
	Mac: number[]
	Leds: number[]

	CurrentValue: number
	AxisX: number[]
	AxisY: number[]
}
