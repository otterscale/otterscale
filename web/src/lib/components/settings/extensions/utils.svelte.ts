import { SvelteURL } from 'svelte/reactivity';

import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import { page } from '$app/state';

import type { ExtensionsBundleType } from './types';

const key = 'bundle';

function getAccordionValue(): ExtensionsBundleType[] | undefined {
	return (page.url.searchParams.getAll(key) as ExtensionsBundleType[]) ?? undefined;
}

function installExtensions(scope: string, extensions: ExtensionsBundleType[]) {
	const basePath = resolve('/(auth)/scope/[scope]/settings/extensions', {
		scope: scope
	});
	const url = new SvelteURL(basePath, window.location.origin);
	extensions.forEach((extension) => url.searchParams.append(key, extension));
	goto(url.pathname + url.search); // eslint-disable-line svelte/no-navigation-without-resolve
}

export { getAccordionValue, installExtensions };
