<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		key
	}: {
		key: User_Key;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {key} />
	</Layout.ActionItem>
</Layout.Actions>
