import type { Component } from 'svelte';

import CronJob from '$lib/components/dynamic-form/cronjob/create-sheet.svelte';

type CreatorProps = { schema?: any };
type CreatorType = Component<CreatorProps> | null;

import Default from './default.svelte';

function getCreator(kind: string): CreatorType {
	switch (kind) {
		case 'CronJob':
			return CronJob as CreatorType;
		default:
			return Default as CreatorType;
	}
}

export { getCreator };
export type { CreatorType };
