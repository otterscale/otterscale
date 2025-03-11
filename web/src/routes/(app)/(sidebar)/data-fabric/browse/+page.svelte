<script lang="ts">
	import { type Node, type Edge } from '@xyflow/svelte';

	import '@xyflow/svelte/dist/style.css';
	import { FabricFlow } from '$lib/components/fabric';
	import { listConnectors, listPipelines } from '$lib/pb';
	import { onMount } from 'svelte';

	const nodes: Node[] = [];
	const edges: Edge[] = [];

	let mounted = false;

	onMount(async () => {
		const filters = [`kind='source'`, `kind='destination'`];
		const connectors = await Promise.all(filters.map((f) => listConnectors(f)));
		connectors.flat().forEach((connector) => {
			nodes.push({
				type: connector.kind,
				id: connector.id,
				data: { ...connector },
				position: { x: 0, y: 0 }
			});
		});

		const pipelines = await listPipelines(`source != '' && destination != ''`);
		pipelines.forEach((pipeline) => {
			edges.push({
				id: pipeline.id,
				source: pipeline.source.id,
				target: pipeline.destination.id
			} as Edge);
		});

		mounted = true;
	});
</script>

{#if mounted}
	<FabricFlow {nodes} {edges} horizontal />
{:else}
	Loading
{/if}
