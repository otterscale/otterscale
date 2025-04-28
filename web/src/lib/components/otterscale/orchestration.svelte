<script lang="ts">
	// External package imports
	import dagre from '@dagrejs/dagre';
	import {
		Background,
		BackgroundVariant,
		Controls,
		MiniMap,
		Panel,
		Position,
		SvelteFlow,
		type ColorModeClass,
		type Edge,
		type Node
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import { writable } from 'svelte/store';
	import { Input } from '$lib/components/ui/input';

	// Internal modules/types
	import type { Instance } from './interfaces';
	import { mode } from 'mode-watcher';

	// Component imports
	import CreateButton from './create.svelte';
	import OrchestrationDrawer from './orchestration-drawer.svelte';
	import OrchestrationNodeMAAS from './orchestration-node-maas.svelte';
	import OrchestrationNodeJUJU from './orchestration-node-juju.svelte';
	import OrchestrationNodeKubernetes from './orchestration-node-kubernetes.svelte';

	const nodeTypes = {
		JUJU: OrchestrationNodeJUJU,
		MAAS: OrchestrationNodeMAAS,
		Kubernetes: OrchestrationNodeKubernetes
	};

	function getLayoutedElements(nodes: Node[], edges: Edge[]) {
		const direction = horizontal ? 'LR' : 'TB';
		const nodeWidth = horizontal ? 100 : 100;
		const nodeHeight = horizontal ? 100 : 100;

		dagreGraph.setGraph({ rankdir: direction });

		nodes.forEach((node) => {
			dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});
		edges.forEach((edge) => {
			dagreGraph.setEdge(edge.source, edge.target);
		});

		dagre.layout(dagreGraph);

		nodes.forEach((node) => {
			const nodeWithPosition = dagreGraph.node(node.id);
			node.position = {
				x: nodeWithPosition.x - nodeWidth / 2,
				y: nodeWithPosition.y - nodeHeight / 2
			};
			node.targetPosition = horizontal ? Position.Left : Position.Top;
			node.sourcePosition = horizontal ? Position.Right : Position.Bottom;
			node.data.horizontal = horizontal;
		});

		return { nodes, edges };
	}

	function reset() {
		$nodeStore.forEach((node) => {
			node.class = undefined;
			$nodeStore = $nodeStore;
		});
		$edgeStore.forEach((edge) => {
			edge.style = undefined;
			$edgeStore = $edgeStore;
		});
	}

	function focus({ id, isEdge, all }: { id: string; isEdge: boolean; all?: boolean }) {
		if (!id) {
			return;
		}
		if (isEdge) {
			$edgeStore.forEach((edge) => {
				if (edge.id === id) {
					edge.style = `stroke-width: 2px; stroke: #f97316;`;
					$edgeStore = $edgeStore;

					focus({ id: edge.source, isEdge: false });
					focus({ id: edge.target, isEdge: false });
				}
			});
			return;
		}
		$nodeStore.forEach((node) => {
			if (node.id === id) {
				node.class = `rounded-full`;
				$nodeStore = $nodeStore;

				if (all) {
					$edgeStore.forEach((edge) => {
						if (edge.source === node.id) {
							focus({ id: edge.id, isEdge: true });
						}
						if (edge.target === node.id) {
							focus({ id: edge.id, isEdge: true });
						}
					});
				}
			}
		});
	}

	let filterPattern = $state('');
	function filter() {
		// Update node visibility based on the filter pattern
		$nodeStore.forEach((node) => {
			node.hidden = !node.id.toLowerCase().includes(filterPattern.toLowerCase());
		});

		// Update edge visibility based on the visibility of connected nodes
		$edgeStore.forEach((edge) => {
			const sourceNode = $nodeStore.find((node) => node.id === edge.source);
			const targetNode = $nodeStore.find((node) => node.id === edge.target);

			// Hide the edge if either the source or target node is hidden
			edge.hidden = sourceNode?.hidden || targetNode?.hidden;
		});

		// Trigger reactivity
		$nodeStore = $nodeStore;
		$edgeStore = $edgeStore;
	}

	let {
		nodes,
		edges,
		horizontal = false
	}: {
		nodes: Node[];
		edges: Edge[];
		horizontal: boolean;
	} = $props();

	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(nodes, edges);
	const nodeStore = writable(layoutedNodes);
	const edgeStore = writable(layoutedEdges);

	let isDrawerOpen = $state(Object.fromEntries(nodes.map((node) => [node.id, false])));
	let isConnected = $state(false);
	let source = $state<Instance>();
	let destination = $state<Instance>();
</script>

{#each nodes as node}
	<drawer class="hidden">
		<OrchestrationDrawer {node} bind:open={isDrawerOpen[node.id]} />
	</drawer>
{/each}

<SvelteFlow
	nodes={nodeStore}
	edges={edgeStore}
	{nodeTypes}
	defaultEdgeOptions={{ animated: true }}
	colorMode={$mode as ColorModeClass}
	proOptions={{ hideAttribution: true }}
	on:nodeclick={(event) => {
		reset();
		focus({ id: event.detail.node.id, isEdge: false, all: true });
		isDrawerOpen[event.detail.node.id] = true;
	}}
	on:edgeclick={(event) => {
		reset();
		focus({ id: event.detail.edge.id, isEdge: true });
	}}
	on:paneclick={() => reset()}
	onconnect={(connection) => {
		const sourceNode = $nodeStore.find((n) => n.id === connection.source);
		const targetNode = $nodeStore.find((n) => n.id === connection.target);
		if (sourceNode && targetNode) {
			source = sourceNode.data as unknown as Instance;
			destination = targetNode.data as unknown as Instance;
			isConnected = true;
		}
	}}
>
	<Background variant={BackgroundVariant.Dots} />
	<Controls orientation="horizontal" />
	<MiniMap pannable zoomable position="bottom-right" />
	<Panel position="top-right" class="px-6">
		<div class="grid justify-items-end gap-3">
			<Input
				placeholder="Filter"
				bind:value={filterPattern}
				oninput={() => {
					filter();
					reset();
				}}
			/>
			<CreateButton />
		</div>
	</Panel>
</SvelteFlow>
