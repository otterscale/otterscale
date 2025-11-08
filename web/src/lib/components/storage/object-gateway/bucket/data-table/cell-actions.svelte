<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		bucket
	}: {
		bucket: Bucket;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {bucket} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {bucket} />
	</Layout.ActionItem>
</Layout.Actions>
