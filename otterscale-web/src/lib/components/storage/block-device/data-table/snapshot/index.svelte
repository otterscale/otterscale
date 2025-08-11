<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Sheet from '$lib/components/ui/sheet';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import { DataTable } from './data-table';

	const selectedFacility = 'ceph-mon';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		image,
		images: data = $bindable()
	}: {
		selectedScopeUuid: string;
		image: Image;
		images: Writable<Image[]>;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{image.snapshots.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw] p-4">
			<DataTable {selectedScopeUuid} {selectedFacility} {image} bind:images={data} />
		</Sheet.Content>
	</Sheet.Root>
</div>
