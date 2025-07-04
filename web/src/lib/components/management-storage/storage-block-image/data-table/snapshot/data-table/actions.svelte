<script lang="ts">
	import type { Image, User_Key } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Protect from './protect.svelte';
	import Rollback from './rollback.svelte';
	import Unprotect from './unprotect.svelte';

	let {
		selectedScope,
		selectedFacility,
		image,
		snapshot,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		image: Image;
		snapshot: User_Key;
		data: Writable<Image[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Rollback {selectedScope} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Protect {selectedScope} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Unprotect {selectedScope} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScope} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
