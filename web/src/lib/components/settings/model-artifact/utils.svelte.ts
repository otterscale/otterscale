import { SvelteURLSearchParams } from 'svelte/reactivity';

import type { HuggingFaceModel } from './types';

async function fetchModels(
	author: string,
	filter: string = '',
	sort: 'downloads' | 'likes' = 'downloads',
	limit: number = 10,
): Promise<HuggingFaceModel[]> {
	const base = 'https://huggingface.co/api/models';
	const queryParameters = new SvelteURLSearchParams({
		author: author,
		filter: filter,
		sort: sort,
		limit: String(limit),
		direction: '-1',
	});

	try {
		const response = await fetch(`${base}?${queryParameters}`);
		if (!response.ok) {
			throw new Error(`Failed to fetch models: ${response.status} ${response.statusText}`);
		}
		const data = await response.json();
		return data as HuggingFaceModel[];
	} catch (error) {
		throw new Error(error instanceof Error ? error.message : String(error));
	}
}

export { fetchModels };
