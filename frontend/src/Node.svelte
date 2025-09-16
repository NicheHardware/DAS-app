<script lang="ts">
	import { onMount } from "svelte"
	import UPlot from "./UPlot.svelte"

	export let axisX: number[] = []
	export let axisY: number[] = []

	export let nodeType: string = "Temperature"
	export let unit: string = "°C"
	export let digit: number = 1
	export let nodeColor: string = "#f73"
	export let battery: number = 30 + Math.random() * 70
	export let RSSI: number = -80 + Math.random() * 50

	let series = [
		{
			label: "Time",
			value: (self, UTC) => (UTC == null ? "" : formatLocalTime(UTC)),
		},
		{
			label: nodeType,
			stroke: nodeColor,
			value: (self, rawValue) =>
				rawValue == null ? "" : rawValue.toFixed(digit) + unit,
		},
	]

	export let currentValue: number =
		15 + Math.random() * 10 + Math.random() * 10 + Math.random() * 10

	export let mac: number[] = Array.from({ length: 6 }, () =>
		Math.floor(Math.random() * 256)
	)
	mac[0] = (mac[0] & 0b11111110) | 0b00000010

	let macAddress: string = mac
		.map((b) => b.toString(16).padStart(2, "0")) // 2-digit hex
		.join(":")
		.toUpperCase()
	let leds: number[] = mac.slice(1).map((b) => b & 1)

	let data: number[][] = [axisX, axisY]
	function formatLocalTime(timestamp: number | Date): string {
		const date = new Date(timestamp)

		// Get local time components
		const hours = date.getHours().toString().padStart(2, "0")
		const minutes = date.getMinutes().toString().padStart(2, "0")
		const seconds = date.getSeconds().toString().padStart(2, "0")

		return `${hours}:${minutes}:${seconds}`
	}
</script>

<div class="flex flex-row gap-1 m-1">
	<div class="flex-0 my-1 border-2 rounded-xl border-gray-400">
		<div class="w-40 mt-2 flex flex-row gap-1 justify-center">
			{#each leds as led}
				{#if led}
					<div class="w-3 h-3 m-1 rounded-full bg-neutral-50" />
				{:else}
					<div
						class="w-3 h-3 m-1 rounded-full border-2 border-neutral-600"
					/>
				{/if}
			{/each}
		</div>
		<div class="w-40 m-1 flex justify-center">
			<span class={`text-2xl font-bold`} style="color: {nodeColor}"
				>{nodeType}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<span class="text-4xl font-bold"
				>{currentValue.toFixed(digit)} {unit}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<div class="flex flex-col text-xs text-neutral-500">
				<span class="">{macAddress}</span>
				<span class="">RSSI:{RSSI.toFixed(0)}dB</span>
				<span class="">BATT:{battery.toFixed(0)}%</span>
			</div>
		</div>
	</div>
	<div class="flex-1 my-1 border-2 rounded-xl border-gray-400">
		<UPlot {series} {data} />
	</div>
</div>
