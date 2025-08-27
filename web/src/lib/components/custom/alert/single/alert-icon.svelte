<script lang="ts" module>
	import { cn } from '$lib/utils';
	import type { AlertVariant } from './alert.svelte';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	export const alertIcons: Map<AlertVariant, string> = new Map([
		['default', 'ph:terminal'],
		['information', 'ph:info'],
		['warning', 'ph:warning'],
		['destructive', 'ph:prohibit']
	]);

	export function getAlertIcon(variant: AlertVariant): string {
		return alertIcons.get(variant) ?? 'ph:terminal';
	}
</script>

<Icon data-slot="alert-icon" icon={getAlertIcon(getContext('variant'))} />
