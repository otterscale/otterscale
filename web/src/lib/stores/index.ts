import { type Writable, writable } from 'svelte/store';

import { resolve } from '$app/paths';
import { m } from '$lib/paraglide/messages';
import type { Path } from '$lib/path';

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

	// Notification
	notifications: Writable<Notification[]>;
}

// Create stores
const createStores = (): AppStores => ({
	breadcrumbs: writable<Path[]>([{ title: m.home(), url: resolve('/') }]),
	// temp
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
export const { breadcrumbs, notifications } = createStores();
