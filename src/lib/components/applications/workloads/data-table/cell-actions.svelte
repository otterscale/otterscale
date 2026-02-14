<script lang="ts" module>
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import type { Application } from '../types';
	import Restart from './action-restart.svelte';
	import Scale from './action-scale.svelte';
</script>

<script lang="ts">
	let {
		application,
		scope,
		reloadManager
	}: { application: Application; scope: string; reloadManager: ReloadManager } = $props();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>
		{m.actions()}
	</Actions.Label>
	<Actions.Separator />
	<Actions.Item>
		<Restart {application} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Scale {application} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
