<script lang="ts" module>
	import Button from '$lib/components/ui/button/button.svelte';
	import dagre from '@dagrejs/dagre';
	import { Background, Controls, MiniMap, Panel, Position, SvelteFlow, type Edge, type Node } from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import FlowEdge from './flow-edge.svelte';
	import FlowNodeGPU from './flow-node-gpu.svelte';
	import FlowNodeMachine from './flow-node-machine.svelte';
	import FlowNodeModel from './flow-node-model.svelte';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts">
	let isHorizontal = $state(false);
	const defaultNodeWidth = 200;
	const nodeHeight = 100;
	function getLayoutedElements(nodes: Node[], edges: Edge[]) {
		dagreGraph.setGraph({ rankdir: isHorizontal ? 'LR' : 'TB' });

		nodes.forEach((node) => {
			dagreGraph.setNode(node.id, { width: defaultNodeWidth, height: nodeHeight });
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
					x: nodeWithPosition.x * 2 + defaultNodeWidth,
					y: nodeWithPosition.y * 1.5 + nodeHeight,
				},
			};
		});

		return { nodes: layoutedNodes, edges };
	}

	function onLayout() {
		isHorizontal = !isHorizontal;
		const layoutedElements = getLayoutedElements(nodes, edges);

		nodes = layoutedElements.nodes;
		edges = layoutedElements.edges;
	}

	function onFocus(nodeId: string) {
		const visited = new Set<string>();
		const toFocus: string[] = [];
		const edgeToFocus: Set<string> = new Set();

		function recurse(currentId: string) {
			if (visited.has(currentId)) return;
			visited.add(currentId);
			toFocus.push(currentId);

			edges.forEach((edge) => {
				if (edge.target === currentId) {
					edgeToFocus.add(edge.id);
					recurse(edge.source);
				}
			});
		}

		recurse(nodeId);

		nodes = nodes.map((node) => ({
			...node,
			selected: toFocus.includes(node.id),
		}));

		edges = edges.map((edge) => ({
			...edge,
			selected: edgeToFocus.has(edge.id),
		}));
	}

	let { initialNodes, initialEdges }: { initialNodes: Node[]; initialEdges: Edge[] } = $props();

	const dagreGraph = new dagre.graphlib.Graph();
	dagreGraph.setDefaultEdgeLabel(() => ({}));

	const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(initialNodes, initialEdges);
	let nodes = $state.raw<Node[]>(layoutedNodes);
	let edges = $state.raw<Edge[]>(layoutedEdges);
</script>

<SvelteFlow
	bind:nodes
	nodeTypes={{
		gpu: FlowNodeGPU,
		machine: FlowNodeMachine,
		model: FlowNodeModel,
	}}
	bind:edges
	edgeTypes={{
		edge: FlowEdge,
	}}
	fitView
	defaultEdgeOptions={{ animated: true }}
	class="min-h-[100vh] border"
	onnodeclick={(e) => {
		onFocus(e.node.id);
	}}
>
	<Background />
	<MiniMap />
	<Controls />
	<Panel position="top-right">
		<Button onclick={() => onLayout()}>
			<Icon icon={isHorizontal ? 'ph:arrow-elbow-up-right-bold' : 'ph:arrow-elbow-left-down-bold'} />
		</Button>
	</Panel>
</SvelteFlow>
