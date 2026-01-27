import type { ComponentType } from 'svelte';

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

function getResourceInspector(resource: string): ComponentType | null {
	if (resource === 'workspaces') {
		return Workspaces;
	}
	return Default;
}

export { getResourceInspector };
