import type { Component } from 'svelte';

type ViewerProps = { object: any };
type ViewerType = Component<ViewerProps>;

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

function getResourceViewer(resource: string): ViewerType {
	if (resource === 'workspaces') {
		return Workspaces as ViewerType;
	}
	return Default as ViewerType;
}

export type { ViewerType };
export { getResourceViewer };
