import type { Component } from 'svelte';

import CronJobActions from './cronjob/actions.svelte';
import CronJobCreate from './cronjob/create.svelte';
import DefaultActions from './default/actions.svelte';
import DefaultCreate from './default/create.svelte';
import WorkspaceActions from './workspace/actions.svelte';
import WorkspaceCreate from './workspace/create.svelte';

type CreateType = Component<{ schema?: any }> | null;
type ActionsType = Component<{
	row?: any;
	schema?: any;
	object?: any;
	cluster?: string;
	onsuccess?: () => void;
}> | null;

function getCreate(kind: string): CreateType {
	switch (kind) {
		case 'CronJob':
			return CronJobCreate as CreateType;
		case 'Workspace':
			return WorkspaceCreate as CreateType;
		default:
			return DefaultCreate as CreateType;
	}
}

function getActions(kind: string): ActionsType {
	switch (kind) {
		case 'CronJob':
			return CronJobActions as ActionsType;
		case 'Workspace':
			return WorkspaceActions as ActionsType;
		default:
			return DefaultActions as ActionsType;
	}
}

export { getActions, getCreate };
export type { ActionsType, CreateType };
