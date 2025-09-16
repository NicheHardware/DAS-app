<script lang="ts">
	import Node from "./Node.svelte"
	import { LogPrintln } from "../wailsjs/go/main/App.js"

	const debugDataCount = 200
	const debugNodeCount = 16
	let timelineDbg: number[] = []
	for (let index = 0; index < debugDataCount; index++) {
		timelineDbg.push(Date.now() + (index - debugDataCount) * 1000)
	}

	let nodes = []
	for (let index = 0; index < debugNodeCount; index++) {
		if (Math.random() > 0.68) {
			nodes.push({
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

				axisY: [],
				battery: 30 + Math.random() * 70,
				RSSI: -80 + Math.random() * 50,
				mac: Array.from({ length: 6 }, () =>
					Math.floor(Math.random() * 256)
				),
			})
		} else {
			nodes.push({
				nodeType: "Temperature",
				unit: "°C",
				digit: 1,
				nodeColor: "#f73",
				currentValue:
					15 +
					Math.random() * 10 +
					Math.random() * 10 +
					Math.random() * 10,

				axisY: [],
				battery: 30 + Math.random() * 70,
				RSSI: -80 + Math.random() * 50,
				mac: Array.from({ length: 6 }, () =>
					Math.floor(Math.random() * 256)
				),
			})
		}
	}
	nodes.forEach((node) => {
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
				for (let index = 0; index < debugDataCount; index++) {
					randP +=
						0.4 * Math.random() - 0.2 - ((randNow - 35) / 20) * 0.2
					randP = Math.max(-2, Math.min(2, randP))
					randNow += randP
					node.axisY.push(randNow + 1 * Math.random())
				}
				break
			case "Voltage":
				randNow =
					0.9 +
					Math.random() * 0.1 +
					Math.random() * 0.1 +
					Math.random() * 0.1 +
					Math.random() * 0.1 +
					Math.random() * 0.1
				for (let index = 0; index < debugDataCount; index++) {
					randP +=
						0.01 * Math.random() -
						0.005 -
						((randNow - 1.15) / 0.25) * 0.005
					randP = Math.max(-0.05, Math.min(0.05, randP))
					randNow += randP
					node.axisY.push(randNow + 0.02 * Math.random())
				}
				break
		}
	})
</script>

<main>
	<h1 class="text-3xl font-bold">DAS</h1>
	{#each nodes as node}
		<Node axisX={timelineDbg} {...node} />
	{/each}
</main>
