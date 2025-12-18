<script lang="ts" module>
	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Log from './action-log.svelte';
</script>

<script lang="ts">
	let {
		pod,
		scope,
		namespace,
		reloadManager
	}: { pod: Application_Pod; scope: string; namespace: string; reloadManager: ReloadManager } =
		$props();

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
		<Log {pod} {scope} {namespace} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {pod} {scope} {namespace} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
