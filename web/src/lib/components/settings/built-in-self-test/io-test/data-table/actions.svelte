<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Retest from './action-retest.svelte';
	import View from './action-view.svelte';
</script>

<script lang="ts">
	let {
		testResult,
		scope,
		reloadManager
	}: {
		testResult: TestResult;
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
		<View {testResult} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Retest {testResult} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {testResult} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
