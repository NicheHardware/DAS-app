<script lang="ts">
	import Node from "./Node.svelte"
	import type { TypeNode } from "./Node"
	import { LogPrintln } from "../wailsjs/go/main/App.js"
	import { randomNode, randomData } from "./Node"

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
	<h1 class="text-3xl font-bold">DAS</h1>
	{#each nodes as node}
		<Node {node} />
	{/each}
</main>
