import type { Component } from 'svelte';

type InspectorProps = { object: any };

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

function getResourceInspector(resource: string): Component<InspectorProps> {
	if (resource === 'workspaces') {
		return Workspaces as Component<InspectorProps>;
	}
	return Default as Component<InspectorProps>;
}

export { getResourceInspector };
