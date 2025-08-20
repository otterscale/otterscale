<script lang="ts">
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import { Layout } from '$lib/components/custom/data-table';
	import Delete from './action-delete.svelte';
	import Retest from './action-retest.svelte';
	import View from './action-view.svelte';

	let {
		testResult
	}: {
		testResult: TestResult;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
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
