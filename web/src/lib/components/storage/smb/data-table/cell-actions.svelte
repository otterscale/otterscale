<script lang="ts" module>
	import type { SMBShare } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		scope,
		facility,
		namespace,
		smbShare,
		reloadManager
	}: {
		scope: string;
		facility: string;
		namespace: string;
		reloadManager: ReloadManager;
		smbShare: SMBShare;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {scope} {facility} {namespace} {smbShare} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
