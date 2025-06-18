<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Bucket } from './types';

	let { bucket }: { bucket: Bucket } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {bucket} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {bucket} />
	</Layout.ActionItem>
</Layout.Actions>
