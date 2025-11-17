<script lang="ts" module>
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import type { LargeLanguageModel } from '../type';
	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		model,
		reloadManager
	}: {
		model: LargeLanguageModel;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {model} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
