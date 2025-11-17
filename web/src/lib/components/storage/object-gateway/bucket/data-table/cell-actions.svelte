<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		bucket,
		scope,
		reloadManager
	}: {
		bucket: Bucket;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {bucket} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {bucket} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
