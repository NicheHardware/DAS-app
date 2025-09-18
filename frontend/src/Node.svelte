<script lang="ts">
	import { onMount } from "svelte"
	import UPlot from "./UPlot.svelte"
	import type { TypeNode } from "./Node"
	import { randomNode } from "./Node"

	export let node: TypeNode = randomNode()

	let series = [
		{
			label: "Time",
			value: (_: null, UTC: number) =>
				UTC == null ? "" : formatLocalTime(UTC),
		},
		{
			label: node.nodeType,
			stroke: node.nodeColor,
			value: (_: null, rawValue: number) =>
				rawValue == null
					? ""
					: rawValue.toFixed(node.digit) + node.unit,
		},
	]

	let macAddress: string = node.mac
		.map((b) => b.toString(16).padStart(2, "0")) // 2-digit hex
		.join(":")
		.toUpperCase()

	let data: number[][]
	$: data = [node.axisX, node.axisY]

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
			{#each node.leds as led}
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
			<span class={`text-2xl font-bold`} style="color: {node.nodeColor}"
				>{node.nodeType}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<span class="text-4xl font-bold"
				>{node.currentValue.toFixed(node.digit)} {node.unit}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<div class="flex flex-col text-xs text-neutral-500">
				<span class="">{macAddress}</span>
				<span class="">RSSI:{node.RSSI.toFixed(0)}dB</span>
				<span class="">BATT:{node.battery.toFixed(0)}%</span>
			</div>
		</div>
	</div>
	<div class="flex-1 my-1 border-2 rounded-xl border-gray-400">
		<UPlot {series} {data} />
	</div>
</div>
