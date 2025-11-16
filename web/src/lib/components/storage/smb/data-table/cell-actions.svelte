<script lang="ts" module>
	import type { SMBShare } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		smbShare,
		scope,
		namespace,
		reloadManager
	}: {
		smbShare: SMBShare;
		scope: string;
		namespace: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {smbShare} {scope} {namespace} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
