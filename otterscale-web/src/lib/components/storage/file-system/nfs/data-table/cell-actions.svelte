<script lang="ts">
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from './actions-delete.svelte';
	import Edit from './actions-edit.svelte';
	import Grant from './actions-grant-export-access.svelte';
	import Revoke from './actions-revoke-export-access.svelte';

	let {
		subvolume
	}: {
		subvolume: Subvolume;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Grant {subvolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Revoke {subvolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Edit {subvolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {subvolume} />
	</Layout.ActionItem>
</Layout.Actions>
