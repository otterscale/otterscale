import type { PluginsBundleType } from './types';

import { goto } from '$app/navigation';
import { page } from '$app/state';
import { dynamicPaths } from '$lib/path';

const key = 'ps';

function getAccordionValue(): PluginsBundleType[] | undefined {
	return (page.url.searchParams.getAll(key) as PluginsBundleType[]) ?? undefined;
}

function installPlugins(plugins: PluginsBundleType[]) {
	goto(
		`${dynamicPaths.settings(page.params.scope).url}/plugins?${plugins.map((plugin) => `${key}=${plugin}`).join('&')}`,
	);
}

export { getAccordionValue, installPlugins };
