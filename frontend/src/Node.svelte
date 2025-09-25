<script lang="ts">
	import { onMount } from "svelte"
	import UPlot from "./UPlot.svelte"
	import type { TypeNode } from "./Node"

	export let node: TypeNode

	export let showChart: boolean = false
	let series = [
		{
			label: "Time",
			value: (_: null, UTC: number) =>
				UTC == null ? "" : formatLocalTime(UTC),
		},
		{
			label: node.Type,
			stroke: node.Color,
			value: (_: null, rawValue: number) =>
				rawValue == null
					? ""
					: rawValue.toFixed(node.Digit) + node.Unit,
		},
	]

	let macAddress: string = node.Mac.map((b) =>
		b.toString(16).padStart(2, "0")
	) // 2-Digit hex
		.join(":")
		.toUpperCase()

	let data: number[][]
	$: data = [node.AxisX, node.AxisY]

	function formatLocalTime(timestamp: number | Date): string {
		const date = new Date(timestamp)

		// Get local time components
		const hours = date.getHours().toString().padStart(2, "0")
		const minutes = date.getMinutes().toString().padStart(2, "0")
		const seconds = date.getSeconds().toString().padStart(2, "0")

		return `${hours}:${minutes}:${seconds}`
	}
</script>

<div class="flex gap-1">
	<div class="flex-0 border-2 rounded-xl border-gray-400">
		<div class="w-40 mt-2 flex flex-row gap-1 justify-center">
			{#each node.Leds as led}
				{#if led}
					<div class="w-3 h-3 m-1 rounded-full bg-neutral-50"></div>
				{:else}
					<div
						class="w-3 h-3 m-1 rounded-full border-2 border-neutral-600"
					></div>
				{/if}
			{/each}
		</div>
		<div class="w-40 m-1 flex justify-center">
			<span class={`text-xl font-bold`} style="color: {node.Color}"
				>{node.Type}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<span class="text-3xl font-bold text-center"
				>{node.CurrentValue.toFixed(node.Digit)} {node.Unit}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center">
			<div class="flex flex-col text-xs text-neutral-500">
				<span class="">{macAddress}</span>
				<span class="">RSSI:{node.RSSI.toFixed(0)}dB</span>
				<span class="">BATT:{node.Battery.toFixed(0)}%</span>
			</div>
		</div>
	</div>
	{#if showChart}
		<div class="flex-1 border-2 rounded-xl border-gray-400">
			<UPlot {series} {data} />
		</div>
	{/if}
</div>
