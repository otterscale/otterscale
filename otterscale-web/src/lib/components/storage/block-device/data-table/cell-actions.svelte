<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		image
	}: {
		image: Image;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {image} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {image} />
	</Layout.ActionItem>
</Layout.Actions>
