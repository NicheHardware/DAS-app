<script lang="ts">
	import Node from "./Node.svelte"
	import type { TypeNode } from "./Node"
	import { LogPrintln } from "../wailsjs/go/main/App.js"
	import { randomNode, randomData } from "./Node"

	let showChart = false
	const debugDataCount = 200
	const debugNodeCount = 16

	let nodes: TypeNode[] = Array.from({ length: debugNodeCount }, () =>
		randomNode()
	)
	nodes.forEach((node) => {
		;[node.axisX, node.axisY] = randomData(node, debugDataCount)
	})
</script>

<main>
	<h1 class="text-3xl font-bold">DAS Console</h1>
	<div class="flex gap-1 flex-wrap" class:flex-col={showChart} class:flex-row={!showChart}>
	{#each nodes as node}
		<Node {node} />
	{/each}
	</div>
</main>
