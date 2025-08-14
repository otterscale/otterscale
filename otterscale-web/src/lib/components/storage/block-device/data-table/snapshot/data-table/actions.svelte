<script lang="ts" module>
	import type { Image, Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Protect from './protect.svelte';
	import Rollback from './rollback.svelte';
	import Unprotect from './unprotect.svelte';

	const selectedFacility = 'ceph-mon';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		image,
		snapshot,
		images: data = $bindable()
	}: {
		selectedScopeUuid: string;
		image: Image;
		snapshot: Image_Snapshot;
		images: Writable<Image[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Rollback {selectedScopeUuid} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={snapshot.protected}>
		<Protect {selectedScopeUuid} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem disabled={!snapshot.protected}>
		<Unprotect {selectedScopeUuid} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} {image} {snapshot} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
