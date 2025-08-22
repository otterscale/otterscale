<script lang="ts" module>
	import type { Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { m } from '$lib/paraglide/messages';
	import Delete from './action-delete.svelte';
	import Protect from './action-protect.svelte';
	import Rollback from './action-rollback.svelte';
	import Unprotect from './action-unprotect.svelte';
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Image_Snapshot;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
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
