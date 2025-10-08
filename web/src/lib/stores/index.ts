import { writable, type Writable } from 'svelte/store';

import { PremiumTier_Level, type PremiumTier } from '$lib/api/environment/v1/environment_pb';
import type { Essential } from '$lib/api/orchestrator/v1/orchestrator_pb';
import type { Scope } from '$lib/api/scope/v1/scope_pb';
import { staticPaths, type Path } from '$lib/path';

// Types
interface BreadcrumbState {
	parents: Path[];
	current: Path;
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
	breadcrumb: Writable<BreadcrumbState>;

	// Premium Tier
	premiumTier: Writable<PremiumTier>;

	// Scope & Essential
	activeScope: Writable<Scope>;
	currentCeph: Writable<Essential | undefined>;
	currentKubernetes: Writable<Essential | undefined>;

	// Bookmark
	bookmarks: Writable<Path[]>;

	// Notification
	notifications: Writable<Notification[]>;
}

// Create stores
const createStores = (): AppStores => ({
	breadcrumb: writable<BreadcrumbState>({ parents: [], current: staticPaths.home }),
	premiumTier: writable<PremiumTier>({ level: PremiumTier_Level.BASIC } as PremiumTier),
	activeScope: writable<Scope>(),
	currentCeph: writable<Essential | undefined>(undefined),
	currentKubernetes: writable<Essential | undefined>(undefined),
	bookmarks: writable<Path[]>([]),
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
			updated: new Date(Date.now() - 86400000),
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
			updated: new Date(Date.now() - 3600000),
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
			updated: new Date(Date.now() - 7200000),
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
			updated: new Date(Date.now() - 1800000),
		},
	]),
});

// Export individual stores
export const { breadcrumb, premiumTier, activeScope, currentCeph, currentKubernetes, bookmarks, notifications } =
	createStores();
