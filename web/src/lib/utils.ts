import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { cubicOut } from "svelte/easing";
import type { TransitionConfig } from "svelte/transition";
import { page } from "$app/state";
import * as m from '$lib/paraglide/messages.js';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

type FlyAndScaleParams = {
	y?: number;
	x?: number;
	start?: number;
	duration?: number;
};

export const flyAndScale = (
	node: Element,
	params: FlyAndScaleParams = { y: -8, x: 0, start: 0.95, duration: 150 }
): TransitionConfig => {
	const style = getComputedStyle(node);
	const transform = style.transform === "none" ? "" : style.transform;

	const scaleConversion = (
		valueA: number,
		scaleA: [number, number],
		scaleB: [number, number]
	) => {
		const [minA, maxA] = scaleA;
		const [minB, maxB] = scaleB;

		const percentage = (valueA - minA) / (maxA - minA);
		const valueB = percentage * (maxB - minB) + minB;

		return valueB;
	};

	const styleToString = (
		style: Record<string, number | string | undefined>
	): string => {
		return Object.keys(style).reduce((str, key) => {
			if (style[key] === undefined) return str;
			return str + `${key}:${style[key]};`;
		}, "");
	};

	return {
		duration: params.duration ?? 200,
		delay: 0,
		css: (t) => {
			const y = scaleConversion(t, [0, 1], [params.y ?? 5, 0]);
			const x = scaleConversion(t, [0, 1], [params.x ?? 0, 0]);
			const scale = scaleConversion(t, [0, 1], [params.start ?? 0.95, 1]);

			return styleToString({
				transform: `${transform} translate3d(${x}px, ${y}px, 0) scale(${scale})`,
				opacity: t
			});
		},
		easing: cubicOut
	};
};

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
	return `${url}?callback=${page.url.pathname}`;
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