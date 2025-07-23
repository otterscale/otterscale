<script lang="ts" module>
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Retest from './retest.svelte';
	import View from './view.svelte';
</script>

<script lang="ts">
	let {
		row,
		data = $bindable()
	}: {
		row: Row<TestResult>;
		data: Writable<TestResult[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
				<View testResult={row.original} />
	</Layout.ActionItem>
	<Layout.ActionItem>
				<Retest testResult={row.original} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
				<Delete testResult={row.original} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
