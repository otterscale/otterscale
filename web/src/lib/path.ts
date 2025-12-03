import type { ResolvedPathname } from '$app/types';

export type Path = {
	url: ResolvedPathname;
	title: string;
};

const PATH_ICON_MAP: Record<string, string> = {
	'/applications': 'ph:compass',
	'/registry': 'ph:table',
	'/compute': 'ph:cpu',
	'/configuration': 'ph:gear',
	'/databases': 'ph:database',
	'/machines': 'ph:computer-tower',
	'/models': 'ph:robot',
	'/networking': 'ph:network',
	'/settings': 'ph:sliders-horizontal',
	'/storage': 'ph:hard-drives'
};

const CEPH_PATH_DISABLED_MAP: Record<string, boolean> = {
	'/compute': true,
	'/storage': true
};

const KUBERNETES_PATH_DISABLED_MAP: Record<string, boolean> = {
	'/applications': true,
	'/registry': true,
	'/compute': true,
	'/databases': true,
	'/models': true,
	'/settings': true,
	'/storage': true
};

function findMatchingPath(url: string, pathMap: Record<string, unknown>): string | undefined {
	return Object.keys(pathMap).find((section) => url.endsWith(section));
}

export function getPathIcon(url: string): string {
	const matchedPath = findMatchingPath(url, PATH_ICON_MAP);
	return matchedPath ? PATH_ICON_MAP[matchedPath] : 'ph:circle-dashed';
}

export function getCephPathDisabled(url: string): boolean {
	const matchedPath = findMatchingPath(url, CEPH_PATH_DISABLED_MAP);
	return matchedPath ? CEPH_PATH_DISABLED_MAP[matchedPath] : false;
}

export function getKubernetesPathDisabled(url: string): boolean {
	const matchedPath = findMatchingPath(url, KUBERNETES_PATH_DISABLED_MAP);
	return matchedPath ? KUBERNETES_PATH_DISABLED_MAP[matchedPath] : false;
}
