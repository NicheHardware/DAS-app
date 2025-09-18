<script lang="ts">
	import { onMount, onDestroy } from "svelte"
	import uPlot from "uplot"
	import "uplot/dist/uPlot.min.css"

	export let series = [{}, {}]
	export let data: number[][] = [[], []] // [x[], y[]]

	let container: HTMLDivElement // parent container <div>
	let chart: uPlot
	let resizeObserver: ResizeObserver

	function updateSize() {
		if (!chart) return
		const rect = container.getBoundingClientRect()
		chart.setSize({ width: rect.width, height: rect.height })
	}

	onMount(() => {
		// initial size from container
		const rect = container.getBoundingClientRect()
		chart = new uPlot(
			{
				width: rect.width,
				height: rect.height,
				series,
				scales: { x: { time: true } },
				axes: [{ show: false }, { show: false }],
			},
			data.map((row) => new Float64Array(row)),
			container
		)

		// watch for parent resize
		resizeObserver = new ResizeObserver(updateSize)
		resizeObserver.observe(container)

		return () => {
			resizeObserver.disconnect()
			chart.destroy()
		}
	})

	// update data when parent passes new arrays
	$: chart && chart.setData(data)
</script>

<!-- parent must have some CSS size for this div -->
<div bind:this={container} class="uplot-wrapper m-2 p-2"></div>

<style>
	.uplot-wrapper {
		/* ensures it grows/shrinks with parent */
		width: 95%;
		height: 60%;
	}
</style>
