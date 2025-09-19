<script lang="ts">
	import Node from "./Node.svelte"
	import type { TypeNode } from "./Node"
	import { LogPrintln } from "../wailsjs/go/main/App.js"
	import { randomNode, randomData } from "./Node"

	let showChart = false
	const debugDataCount = 200
	const debugNodeCount = 16
	const headerStyle = "top-0 left-0 right-0 h-14"
	const footerStyle = "bottom-0 left-0 right-0 h-10"

	let nodes: TypeNode[] = Array.from({ length: debugNodeCount }, () =>
		randomNode()
	)
	nodes.forEach((node) => {
		;[node.axisX, node.axisY] = randomData(node, debugDataCount)
	})
</script>

<header
	class={`fixed ${headerStyle} backdrop-blur-sm bg-black/30 border-b z-50 flex items-center px-4`}
>
	<div class="text-3xl font-bold">DAS Console</div>

	<div class="ml-auto flex items-center gap-2">
		<!-- Grid view button (active when showChart is false) -->
		<button
			class="btn btn-outline btn-info btn-sm"
			aria-pressed={!showChart}
			on:click={() => (showChart = false)}
			title="Grid view"
			class:btn-active={!showChart}
		>
			<!-- grid icon -->
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="w-5 h-5"
				viewBox="0 0 24 24"
				fill="currentColor"
				aria-hidden="true"
			>
				<path
					d="M3 3h8v8H3V3zm10 0h8v8h-8V3zM3 13h8v8H3v-8zm10 10V13h8v10h-8z"
				/>
			</svg>
			<span class="sr-only">Grid</span>
		</button>

		<!-- List / chart view button (active when showChart is true) -->
		<button
			class="btn btn-outline btn-info btn-sm"
			aria-pressed={showChart}
			on:click={() => (showChart = true)}
			title="List / Chart view"
			class:btn-active={showChart}
		>
			<!-- list icon -->
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="w-5 h-5"
				viewBox="0 0 24 24"
				fill="currentColor"
				aria-hidden="true"
			>
				<path d="M4 6h16v2H4V6zm0 5h16v2H4v-2zm0 5h16v2H4v-2z" />
			</svg>
			<span class="sr-only">List</span>
		</button>
	</div>
</header>
<div class={headerStyle}></div>
<main class="my-2">
	<div
		class="flex gap-1 flex-wrap"
		class:flex-col={showChart}
		class:flex-row={!showChart}
	>
		{#each nodes as node}
			<Node {node} {showChart} />
		{/each}
	</div>
</main>
<div class={footerStyle}></div>
<footer
	class={`fixed ${footerStyle} border-t z-40 backdrop-blur-sm bg-black/30 flex items-center px-4 gap-2`}
>
	<button class="btn btn-sm">Action 1</button>
	<button class="btn btn-sm">Action 2</button>
	<button class="btn btn-sm">Action 3</button>
	<button class="btn btn-sm">Action 4</button>
	<div class="ml-auto text-sm text-gray-500">Placeholder footer</div>
</footer>

<style>
	:global(html) {
		/* reserve the vertical scrollbar space so layout doesn't shift when content overflows */
		overflow-y: scroll;
	}
</style>
