<script lang="ts" module>
	import { type Manifest } from '$lib/api/registry/v1/registry_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		manifest,
		scope,
		reloadManager
	}: {
		manifest: Manifest;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Delete {manifest} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
