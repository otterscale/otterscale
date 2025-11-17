<script lang="ts" module>
	import type { Image, Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Protect from './action-protect.svelte';
	import Rollback from './action-rollback.svelte';
	import Unprotect from './action-unprotect.svelte';
</script>

<script lang="ts">
	let {
		snapshot,
		image,
		scope,
		reloadManager
	}: {
		snapshot: Image_Snapshot;
		image: Image;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Rollback {snapshot} {image} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={snapshot.protected}>
		<Protect {snapshot} {image} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={!snapshot.protected}>
		<Unprotect {snapshot} {image} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {snapshot} {image} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
