import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { cubicOut } from "svelte/easing";
import type { TransitionConfig } from "svelte/transition";
import { page } from "$app/state";
import * as m from '$lib/paraglide/messages.js';
import { i18n } from "./i18n";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function getFeatureTitle(path: string): string {
	switch (path) {
		case '/tutorial':
			return m.nav_tutorial();
		case '/data-fabric':
			return m.nav_data_fabric();
		case '/explore':
			return m.nav_explore();
		case '/dashboard':
			return m.nav_dashboard();
		case '/applications':
			return m.nav_applications();
		case '/integrations':
			return m.nav_integrations();
		case '/dev-tools':
			return m.nav_dev_tools();
		default:
			return '';
	}
}

export function setCallback(url: string): string {
	return `${url}?callback=${i18n.route(page.url.pathname)}`;
}

export function getCallback(): string {
	const callbackParam = page.url.searchParams.get('callback');
	if (callbackParam) {
		return callbackParam;
	}
	return '/';
}

export function appendCallback(url: string): string {
	const callbackParam = page.url.searchParams.get('callback');
	if (callbackParam) {
		return `${url}?callback=${callbackParam}`;
	}
	return url;
}

const DIVISIONS = [
	{ amount: 60, name: 'seconds' },
	{ amount: 60, name: 'minutes' },
	{ amount: 24, name: 'hours' },
	{ amount: 7, name: 'days' },
	{ amount: 4.34524, name: 'weeks' },
	{ amount: 12, name: 'months' },
	{ amount: Number.POSITIVE_INFINITY, name: 'years' }
] as const;

const formatter = new Intl.RelativeTimeFormat(undefined, {
	numeric: 'auto'
});

export function formatTimeAgo(date: Date) {
	let duration = (new Date(date).getTime() - new Date().getTime()) / 1000;

	for (let i = 0; i <= DIVISIONS.length; i++) {
		const division = DIVISIONS[i];
		if (Math.abs(duration) < division.amount) {
			return formatter.format(Math.round(duration), division.name);
		}
		duration /= division.amount;
	}
}