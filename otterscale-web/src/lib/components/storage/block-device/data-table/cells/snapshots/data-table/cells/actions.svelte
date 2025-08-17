<script lang="ts" module>
	import type { Image, Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from '../actions/delete.svelte';
	import Protect from '../actions/protect.svelte';
	import Rollback from '../actions/rollback.svelte';
	import Unprotect from '../actions/unprotect.svelte';
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Image_Snapshot;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Rollback {snapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={snapshot.protected}>
		<Protect {snapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={!snapshot.protected}>
		<Unprotect {snapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {snapshot} />
	</Layout.ActionItem>
</Layout.Actions>
