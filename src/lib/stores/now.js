import { writable } from 'svelte/store';

function createNowStore() {
	const { subscribe, set } = writable(Date.now());

	setInterval(() => {
		set(Date.now());
	}, 1 * 1000);

	return { subscribe };
}

export const now = createNowStore();
