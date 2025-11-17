<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Retest from './action-retest.svelte';
	import View from './action-view.svelte';
</script>

<script lang="ts">
	let {
		testResult,
		reloadManager
	}: {
		testResult: TestResult;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<View {testResult} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Retest {testResult} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {testResult} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
