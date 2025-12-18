<script lang="ts" module>
	import type { Model } from '$lib/api/model/v1/model_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';

	import Delete from './action-delete.svelte';
	import Update from './action-update.svelte';
</script>

<script lang="ts">
	let {
		model,
		scope,
		reloadManager
	}: {
		model: Model;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

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
		<Update {model} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {model} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
