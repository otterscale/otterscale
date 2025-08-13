<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		bucket,
		buckets: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		bucket: Bucket;
		buckets: Writable<Bucket[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScopeUuid} {selectedFacility} {bucket} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} {bucket} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
