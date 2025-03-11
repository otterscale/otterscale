<script lang="ts">
	import dagre from '@dagrejs/dagre';
	import { writable } from 'svelte/store';
	import {
		SvelteFlow,
		Background,
		Controls,
		BackgroundVariant,
		type ColorModeClass,
		MiniMap,
		type Node,
		type Edge,
		Position
	} from '@xyflow/svelte';

	import '@xyflow/svelte/dist/style.css';
	import { mode } from 'mode-watcher';
	import ConnectorNode from './connector-node.svelte';

	let {
		nodes,
		edges,
		horizontal = false
	}: {
		nodes: Node[];
		edges: Edge[];
		horizontal: boolean;
	} = $props();

	const nodeTypes = {
		source: ConnectorNode,
		destination: ConnectorNode
	};

	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(nodes, edges);

	const nodeStore = writable(layoutedNodes);
	const edgeStore = writable(layoutedEdges);

	function getLayoutedElements(nodes: Node[], edges: Edge[]) {
		const direction = horizontal ? 'LR' : 'TB';
		const nodeWidth = horizontal ? 600 : 450;
		const nodeHeight = horizontal ? 120 : 240;

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

					// find nodes
					focus({ id: edge.source, isEdge: false });
					focus({ id: edge.target, isEdge: false });
				}
			});
			return;
		}
		$nodeStore.forEach((node) => {
			if (node.id === id) {
				node.class = `bg-accent`;
				$nodeStore = $nodeStore;

				// find edges
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
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
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
		}}
		on:edgeclick={(event) => {
			reset();
			focus({ id: event.detail.edge.id, isEdge: true });
		}}
		on:paneclick={() => reset()}
	>
		<Background variant={BackgroundVariant.Dots} />
		<Controls orientation="horizontal" />
		<MiniMap pannable zoomable position="top-right" />
	</SvelteFlow>
</main>
