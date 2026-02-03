import { type Writable, writable } from 'svelte/store';

import { browser } from '$app/environment';
import { resolve } from '$app/paths';
import { m } from '$lib/paraglide/messages';
import type { Path } from '$lib/path';

// Create a writable store that persists to localStorage
function persistentWritable<T>(key: string, initialValue: T): Writable<T> {
	// Get initial value from localStorage if in browser
	const getStoredValue = () => {
		if (!browser) return initialValue;
		try {
			const stored = localStorage.getItem(key);
			return stored ? JSON.parse(stored) : initialValue;
		} catch {
			return initialValue;
		}
	};

	const store = writable<T>(getStoredValue());

	// Subscribe to store changes and persist to localStorage
	if (browser) {
		store.subscribe((value) => {
			try {
				localStorage.setItem(key, JSON.stringify(value));
			} catch (error) {
				console.warn(`Failed to persist ${key} to localStorage:`, error);
			}
		});
	}

	return store;
}

// temp
export interface Notification {
	id: string;
	from: string;
	title: string;
	content: string;
	read: boolean;
	archived: boolean;
	deleted: boolean;
	created: Date;
	updated: Date;
}

interface AppStores {
	// Navigation
	breadcrumbs: Writable<Path[]>;

	// Workspace
	activeWorkspaceName: Writable<string>;

	// Namespace
	activeNamespace: Writable<string>;

	// Notification
	notifications: Writable<Notification[]>;
}

// Create stores
const createStores = (): AppStores => ({
	breadcrumbs: writable<Path[]>([{ title: m.home(), url: resolve('/') }]),
	// Persistent workspace store
	activeWorkspaceName: persistentWritable<string>('otterscale:activeWorkspace', ''),
	// Persistent namespace store
	activeNamespace: persistentWritable<string>('otterscale:activeNamespace', ''),
	notifications: writable<Notification[]>([
		{
			id: '1',
			from: 'system@otterscale.com',
			title: 'Welcome to OtterScale',
			content: 'Your account has been successfully created.',
			read: false,
			archived: false,
			deleted: false,
			created: new Date(Date.now() - 86400000),
			updated: new Date(Date.now() - 86400000)
		}
	])
});

// Export individual stores
export const { breadcrumbs, activeWorkspaceName, activeNamespace, notifications } = createStores();
