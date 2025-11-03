import type { ExtensionsBundleType } from './types';

import { goto } from '$app/navigation';
import { page } from '$app/state';
import { dynamicPaths } from '$lib/path';

const key = 'ps';

function getAccordionValue(): ExtensionsBundleType[] | undefined {
	return (page.url.searchParams.getAll(key) as ExtensionsBundleType[]) ?? undefined;
}

function installExtensions(extensions: ExtensionsBundleType[]) {
	goto(
		`${dynamicPaths.settingsExtensions(page.params.scope).url}?${extensions.map((extension) => `${key}=${extension}`).join('&')}`,
	);
}

export { getAccordionValue, installExtensions };
