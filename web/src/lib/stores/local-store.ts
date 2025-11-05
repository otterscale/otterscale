import { writable, type Writable } from 'svelte/store';

import { browser } from '$app/environment';

export function localStore<T>(key: string, initialValue: T): Writable<T> {
	let storedValue: T = initialValue;

	if (browser) {
		const json = localStorage.getItem(key);
		if (json) {
			try {
				storedValue = JSON.parse(json) as T;
			} catch (e) {
				console.error(`[localStore] Error parsing localStorage key "${key}":`, e);
			}
		}
	}

	const store: Writable<T> = writable(storedValue);

	if (browser) {
		store.subscribe((value) => {
			try {
				if (value === undefined) {
					localStorage.removeItem(key);
					return;
				}
				const valueToStore = JSON.stringify(value);
				localStorage.setItem(key, valueToStore);
			} catch (e) {
				console.error(`[localStore] Could not save to localStorage key "${key}"`, e);
			}
		});

		window.addEventListener('storage', (event: StorageEvent) => {
			if (event.key === key && event.newValue) {
				try {
					const newValue = JSON.parse(event.newValue) as T;
					store.set(newValue);
				} catch (e) {
					console.error(`[localStore] Error parsing storage event value for key "${key}"`, e);
				}
			}
		});
	}

	return store;
}
