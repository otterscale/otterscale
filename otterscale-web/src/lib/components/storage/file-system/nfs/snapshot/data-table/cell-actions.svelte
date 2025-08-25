<script lang="ts" module>
	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		snapshot
	}: {
		snapshot: Subvolume_Snapshot;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {snapshot} />
	</Layout.ActionItem>
</Layout.Actions>
