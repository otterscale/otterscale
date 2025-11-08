<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Update from './action-update.svelte';
</script>

<script lang="ts">
	let {
		row
	}: {
		row: Row<Network_IPRange>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Update ipRange={row.original} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete ipRange={row.original} />
	</Layout.ActionItem>
</Layout.Actions>
