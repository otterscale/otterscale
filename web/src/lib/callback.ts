import { page } from "$app/state";
import { i18n } from "./i18n";

export function setCallback(to: string, from: string): string {
	return `${i18n.resolveRoute(to)}?callback=${i18n.route(from)}`;
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
