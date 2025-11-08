<script lang="ts" module>
	import dagre from '@dagrejs/dagre';
	import {
		Controls,
		MiniMap,
		Panel,
		Position,
		SvelteFlow,
		type Edge,
		type Node
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';

	import { traverse } from '../utils.svelte';

	import FlowEdge from './flow-edge.svelte';
	import FlowNodeGPU from './flow-node-gpu.svelte';
	import FlowNodeMachine from './flow-node-machine.svelte';
	import FlowNodeModel from './flow-node-model.svelte';
	import FlowNodeVGPU from './flow-node-vgpu.svelte';

	import { cn } from '$lib/utils';

	const defaultNodeWidth = 300;
	const defaultNodeHeight = 150;
</script>

<script lang="ts">
	let isHorizontal = $state(false);
	function getLayoutedElements(nodes: Node[], edges: Edge[]) {
		dagreGraph.setGraph({ rankdir: isHorizontal ? 'LR' : 'TB' });

		nodes.forEach((node) => {
			dagreGraph.setNode(node.id, { width: defaultNodeWidth, height: defaultNodeHeight });
		});

		edges.forEach((edge) => {
			dagreGraph.setEdge(edge.source, edge.target);
		});

		dagre.layout(dagreGraph);

		const layoutedNodes = nodes.map((node) => {
			const nodeWithPosition = dagreGraph.node(node.id);
			node.targetPosition = isHorizontal ? Position.Left : Position.Top;
			node.sourcePosition = isHorizontal ? Position.Right : Position.Bottom;

			return {
				...node,
				position: {
					x: nodeWithPosition.x - 0.5 * defaultNodeWidth,
					y: nodeWithPosition.y * 1.5 - 0.5 * defaultNodeHeight
				}
			};
		});

		return { nodes: layoutedNodes, edges };
	}

	let {
		initialNodes,
		initialEdges,
		class: className
	}: { initialNodes: Node[]; initialEdges: Edge[]; class?: string } = $props();

	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
		initialNodes,
		initialEdges
	);

	let nodes = $state.raw<Node[]>(layoutedNodes);
	let edges = $state.raw<Edge[]>(layoutedEdges);
</script>

<SvelteFlow
	bind:nodes
	nodeTypes={{
		gpu: FlowNodeGPU,
		vgpu: FlowNodeVGPU,
		machine: FlowNodeMachine,
		model: FlowNodeModel
	}}
	bind:edges
	edgeTypes={{
		edge: FlowEdge
	}}
	fitView
	proOptions={{ hideAttribution: true }}
	class={cn('h-full w-full', className)}
	onnodeclick={(e) => {
		const { nodestoFocus, edgesToFocus } = traverse(edges, e.node.id);

		nodes = nodes.map((node) => ({
			...node,
			selected: nodestoFocus.has(node.id)
		}));

		edges = edges.map((edge) => ({
			...edge,
			selected: edgesToFocus.has(edge.id)
		}));
	}}
>
	<Panel position="top-right" />
	<MiniMap />
	<Controls />
</SvelteFlow>
