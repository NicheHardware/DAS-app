<script lang="ts">
	import UPlot from "./UPlot.svelte"
	import type { TypeNode } from "./Node"

	export let node: TypeNode

	export let showChart: boolean = false
	let selectedMetricIndex = 0
	$: metrics = node.Metrics ?? []
	$: isSingleMetric = metrics.length <= 1
	$: if (selectedMetricIndex >= metrics.length) selectedMetricIndex = 0
	$: selectedMetric = metrics[selectedMetricIndex] ?? metrics[0]
	$: chartKey = selectedMetric
		? `${selectedMetric.Key}-${selectedMetric.Color}`
		: "empty"
	$: series = selectedMetric
		? [
				{
					label: "Time",
					value: (_: null, UTC: number) =>
						UTC == null ? "" : formatLocalTime(UTC),
				},
				{
					label: selectedMetric.Name,
					stroke: selectedMetric.Color,
					value: (_: null, rawValue: number) =>
						rawValue == null
							? ""
							: rawValue.toFixed(selectedMetric.Digit) + selectedMetric.Unit,
				},
		  ]
		: [{}, {}]

	let macAddress: string = node.Mac.map((b) =>
		b.toString(16).padStart(2, "0")
	) // 2-Digit hex
		.join(":")
		.toUpperCase()

	let data: number[][]
	$: data = selectedMetric ? [node.AxisX, selectedMetric.AxisY] : [[], []]

	function formatLocalTime(timestamp: number | Date): string {
		const date = new Date(timestamp)

		// Get local time components
		const hours = date.getHours().toString().padStart(2, "0")
		const minutes = date.getMinutes().toString().padStart(2, "0")
		const seconds = date.getSeconds().toString().padStart(2, "0")

		return `${hours}:${minutes}:${seconds}`
	}

	function selectMetric(index: number) {
		selectedMetricIndex = index
	}

	function getMetricLabel(name: string) {
		return name.slice(0, 3)
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
			<span class={`text-xl font-bold text-center`}>{node.Name}</span
			>
		</div>
		<div class="w-40 m-1 flex justify-center text-sm text-neutral-400">
			<span>{node.TypeName}</span>
		</div>
		<div class="w-40 m-1 px-2 flex flex-col gap-1">
			{#each metrics as metric, index}
				{@const isSelected = selectedMetricIndex === index}
				<button
					class="flex items-center justify-between rounded-md border px-2 py-1 text-left transition-colors"
					class:bg-neutral-800={isSelected}
					class:border-neutral-700={!isSelected}
					on:click={() => selectMetric(index)}
					style:border-color={isSelected ? metric.Color : undefined}
					style:box-shadow={isSelected
						? `0 0 0 1px ${metric.Color} inset`
						: undefined}
				>
					{#if isSingleMetric}
						<span class="w-full text-center font-bold"
							>{metric.CurrentValue.toFixed(metric.Digit)} {metric.Unit}</span
						>
					{:else}
						<span style="color: {metric.Color}">{getMetricLabel(metric.Name)}</span>
						<span class="font-bold"
							>{metric.CurrentValue.toFixed(metric.Digit)} {metric.Unit}</span
						>
					{/if}
				</button>
			{/each}
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
			{#key chartKey}
				<UPlot {series} {data} />
			{/key}
		</div>
	{/if}
</div>
