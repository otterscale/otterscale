import { getIcon } from '@iconify/svelte';

function fuzzLogosIcon(name: string, defaultName: string): string {
	if (!name) {
		return defaultName;
	}
	let icon = `logos:${name}-icon`;
	if (getIcon(icon)) {
		return icon;
	}
	icon = `logos:${name}`;
	if (getIcon(icon)) {
		return icon;
	}
	return defaultName;
}

export { fuzzLogosIcon };
