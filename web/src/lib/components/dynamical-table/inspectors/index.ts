import type { Component } from 'svelte';

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

function getResourceInspector(resource: string): Component | null {
	if (resource === 'workspaces') {
		return Workspaces;
	}
	return Default;
}

export { getResourceInspector };
