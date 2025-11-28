import { SvelteURLSearchParams } from 'svelte/reactivity';

import type { HuggingFaceModel, ModelTag, ModelTagCategory, SortType } from './types';

async function fetchModels(
	author: string,
	tags: string[],
	sort: SortType,
	limit?: number
): Promise<HuggingFaceModel[]> {
	const base = 'https://huggingface.co/api/models';
	const queryParameters = new SvelteURLSearchParams({
		author: author,
		sort: sort,
		direction: '-1'
	});
	tags.forEach((tag) => {
		queryParameters.append('filter', tag);
	});
	if (limit) {
		queryParameters.append('limit', String(limit));
	}

	try {
		const response = await fetch(`${base}?${queryParameters}`);
		if (!response.ok) {
			throw new Error(`Failed to fetch models: ${response.status} ${response.statusText}`);
		}
		const data = await response.json();
		console.log(`${base}?${queryParameters}`);
		console.log(data);
		return data;
	} catch (error) {
		throw new Error(error instanceof Error ? error.message : String(error));
	}
}

async function fetchModelTypes(modelTagCategory: ModelTagCategory): Promise<ModelTag[]> {
	const base = 'https://huggingface.co/api/models-tags-by-type';
	const queryParameters = new SvelteURLSearchParams({
		type: modelTagCategory
	});

	try {
		const response = await fetch(`${base}?${queryParameters}`);
		if (!response.ok) {
			throw new Error(`Failed to fetch model types: ${response.status} ${response.statusText}`);
		}
		const data = await response.json();
		return data[modelTagCategory];
	} catch (error) {
		throw new Error(error instanceof Error ? error.message : String(error));
	}
}

export { fetchModels, fetchModelTypes };
