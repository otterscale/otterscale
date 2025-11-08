<script lang="ts" module>
	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Retest from './action-retest.svelte';
	import View from './action-view.svelte';
</script>

<script lang="ts">
	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<View {testResult} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Retest {testResult} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {testResult} />
	</Layout.ActionItem>
</Layout.Actions>
