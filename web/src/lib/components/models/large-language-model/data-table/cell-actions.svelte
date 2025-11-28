<script lang="ts" module>
	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Update from './action-update.svelte';
</script>

<script lang="ts">
	let {
		model,
		scope,
		reloadManager
	}: {
		model: Model;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Update {model} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {model} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
