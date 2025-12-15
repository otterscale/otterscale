import { type Writable, writable } from 'svelte/store';

import { browser } from '$app/environment';
import { resolve } from '$app/paths';
import { type PremiumTier, PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
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

	// Premium Tier
	premiumTier: Writable<PremiumTier>;

	// Bookmark
	bookmarks: Writable<Path[]>;

	// Notification
	notifications: Writable<Notification[]>;
}

// Helper for persisted store
const createPersistedStore = <T>(key: string, initialValue: T): Writable<T> => {
	if (!browser) return writable(initialValue);

	let data = initialValue;
	try {
		const storedValue = localStorage.getItem(key);
		if (storedValue) {
			data = JSON.parse(storedValue);
		}
	} catch (e) {
		console.warn(`Failed to load ${key} from localStorage`, e);
	}

	const store = writable(data);

	store.subscribe((value) => {
		try {
			localStorage.setItem(key, JSON.stringify(value));
		} catch (e) {
			console.warn(`Failed to save ${key} to localStorage`, e);
		}
	});

	return store;
};

// Create stores
const createStores = (): AppStores => ({
	breadcrumbs: writable<Path[]>([{ title: m.home(), url: resolve('/') }]),
	premiumTier: writable<PremiumTier>({ level: PremiumTier_Level.COMMUNITY } as PremiumTier),
	bookmarks: createPersistedStore<Path[]>('otterscale:bookmarks', []),
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
		},
		{
			id: '2',
			from: 'alerts@otterscale.com',
			title: 'High CPU Usage Alert',
			content: 'Your Kubernetes cluster is experiencing high CPU usage.',
			read: true,
			archived: false,
			deleted: false,
			created: new Date(Date.now() - 3600000),
			updated: new Date(Date.now() - 3600000)
		},
		{
			id: '3',
			from: 'billing@otterscale.com',
			title: 'Monthly Invoice Available',
			content: 'Your monthly invoice for December is now available.',
			read: false,
			archived: false,
			deleted: false,
			created: new Date(Date.now() - 7200000),
			updated: new Date(Date.now() - 7200000)
		},
		{
			id: '4',
			from: 'security@otterscale.com',
			title: 'Security Vulnerability Report',
			content:
				'We have detected several security vulnerabilities in your infrastructure that require immediate attention. Our automated security scanning system has identified the following issues: 1) Outdated Kubernetes version (v1.18.20) with known CVE-2021-25737 vulnerability that allows privilege escalation, 2) Exposed Redis instances on ports 6379 and 6380 without authentication, 3) SSL certificates for *.example.com expiring in 7 days, 4) Unencrypted data transmission detected between microservices in the production namespace, 5) Weak password policies in effect allowing passwords shorter than 12 characters. We strongly recommend immediate action to patch these vulnerabilities. Please review the detailed security report in your dashboard and implement the suggested fixes within 24 hours to maintain compliance with SOC 2 Type II standards.',
			read: false,
			archived: false,
			deleted: false,
			created: new Date(Date.now() - 1800000),
			updated: new Date(Date.now() - 1800000)
		}
	])
});

// Export individual stores
export const { breadcrumbs, premiumTier, bookmarks, notifications } = createStores();
