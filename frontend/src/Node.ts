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

export function randomData(node: TypeNode, count: number) {
	let timeline: number[] = []
	let data: number[] = []
	for (let index = 0; index < count; index++) {
		timeline.push(Date.now() + (index - count) * 1000)
	}

	let randP = 0
	let randNow = 0
	switch (node.Type) {
		case "Temperature":
			randNow =
				15 +
				10 * Math.random() +
				10 * Math.random() +
				10 * Math.random() +
				10 * Math.random()
			for (let index = 0; index < count; index++) {
				randP += 0.4 * Math.random() - 0.2 - ((randNow - 35) / 20) * 0.2
				randP = Math.max(-2, Math.min(2, randP))
				randNow += randP
				data.push(randNow + 1 * Math.random())
			}
		case "Voltage":
			randNow =
				0.9 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1
			for (let index = 0; index < count; index++) {
				randP +=
					0.01 * Math.random() -
					0.005 -
					((randNow - 1.15) / 0.25) * 0.005
				randP = Math.max(-0.05, Math.min(0.05, randP))
				randNow += randP
				data.push(randNow + 0.02 * Math.random())
			}
	}
	return [timeline, data]
}

export function randomNode() {
	let mac: number[] = Array.from({ length: 6 }, () =>
		Math.floor(Math.random() * 256)
	)
	mac[0] = (mac[0] & 0b11111110) | 0b00000010
	let leds = mac.slice(1).map((b) => b & 1)
	if (Math.random() > 0.68) {
		return {
			Type: "Voltage",
			Unit: "V",
			Digit: 2,
			Color: "#37f",
			CurrentValue:
				0.9 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1,

			AxisX: [],
			AxisY: [],
			Battery: 30 + Math.random() * 70,
			RSSI: -80 + Math.random() * 50,
			Mac: mac,
			Leds: leds,
		}
	} else {
		return {
			Type: "Temperature",
			Unit: "°C",
			Digit: 1,
			Color: "#f73",
			CurrentValue:
				15 +
				Math.random() * 10 +
				Math.random() * 10 +
				Math.random() * 10,

			AxisX: [],
			AxisY: [],
			Battery: 30 + Math.random() * 70,
			RSSI: -80 + Math.random() * 50,
			Mac: mac,
			Leds: leds,
		}
	}
}
