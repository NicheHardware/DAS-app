export interface TypeNode {
	nodeType: string
	unit: string
	digit: number
	nodeColor: string
	battery: number
	RSSI: number
	mac: number[]
	leds: number[]

	currentValue: number
	axisX: number[]
	axisY: number[]
}

export function randomData(node: TypeNode, count: number) {
	let timeline: number[] = []
	let data: number[] = []
	for (let index = 0; index < count; index++) {
		timeline.push(Date.now() + (index - count) * 1000)
	}

	let randP = 0
	let randNow = 0
	switch (node.nodeType) {
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
			nodeType: "Voltage",
			unit: "V",
			digit: 2,
			nodeColor: "#37f",
			currentValue:
				0.9 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1 +
				Math.random() * 0.1,

			axisX: [],
			axisY: [],
			battery: 30 + Math.random() * 70,
			RSSI: -80 + Math.random() * 50,
			mac: mac,
			leds: leds,
		}
	} else {
		return {
			nodeType: "Temperature",
			unit: "°C",
			digit: 1,
			nodeColor: "#f73",
			currentValue:
				15 +
				Math.random() * 10 +
				Math.random() * 10 +
				Math.random() * 10,

			axisX: [],
			axisY: [],
			battery: 30 + Math.random() * 70,
			RSSI: -80 + Math.random() * 50,
			mac: mac,
			leds: leds,
		}
	}
}
