import type { Component } from 'svelte';

import CronJob from './cronjob.svelte';
import Default from './default.svelte';
import SimpleApp from './simpleapp.svelte';

type ActionsProps = { row: any; schema?: any; object?: any; onsuccess?: () => void };
type ActionsType = Component<ActionsProps> | null;

function getActions(kind: string): ActionsType {
	switch (kind) {
		case 'CronJob':
			return CronJob as ActionsType;
		case 'SimpleApp':
			return SimpleApp as ActionsType;
		default:
			return Default as ActionsType;
	}
}

export { getActions };
export type { ActionsType };
