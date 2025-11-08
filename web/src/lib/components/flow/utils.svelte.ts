import { type Edge } from '@xyflow/svelte';
import '@xyflow/svelte/dist/style.css';
import { SvelteSet } from 'svelte/reactivity';

function traverse(
	edges: Edge[],
	rootId: string
): { nodestoFocus: SvelteSet<string>; edgesToFocus: SvelteSet<string> } {
	function traverseParents(currentId: string) {
		if (visitedNodes.has(currentId)) return;
		visitedNodes.add(currentId);

		nodestoFocus.add(currentId);

		edges.forEach((edge) => {
			if (edge.target === currentId) {
				edgesToFocus.add(edge.id);
				traverseParents(edge.source);
			}
		});
	}

	function traverseChilds(currentId: string) {
		if (visitedNodes.has(currentId)) return;
		visitedNodes.add(currentId);

		nodestoFocus.add(currentId);

		edges.forEach((edge) => {
			if (edge.source === currentId) {
				edgesToFocus.add(edge.id);
				traverseChilds(edge.target);
			}
		});
	}

	const visitedNodes: SvelteSet<string> = new SvelteSet();
	const nodestoFocus: SvelteSet<string> = new SvelteSet();
	const edgesToFocus: SvelteSet<string> = new SvelteSet();

	visitedNodes.clear();
	traverseParents(rootId);
	visitedNodes.clear();
	traverseChilds(rootId);

	return { nodestoFocus, edgesToFocus };
}

export { traverse };
