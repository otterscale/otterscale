<script lang="ts" module>
	import type { LargeLanguageModel } from '../type';

	import Delete from './action-delete.svelte';
	import Update from './action-update.svelte';

	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		model,
	}: {
		model: LargeLanguageModel;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Update {model} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {model} />
	</Layout.ActionItem>
</Layout.Actions>
