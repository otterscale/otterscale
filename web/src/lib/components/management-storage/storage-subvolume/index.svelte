<script lang="ts" module>
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let selectedScope = $state('b62d195e-3905-4960-85ee-7673f71eb21e');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('');
	let selectedSubvolumeGroup = $state('');
</script>

<main class="space-y-4">
	<Pickers
		bind:selectedScope
		bind:selectedFacility
		bind:selectedVolume
		bind:selectedSubvolumeGroup
	/>

	{#await storageClient.listSubvolumes( { scopeUuid: selectedScope, facilityName: selectedFacility, volumeName: selectedVolume, groupName: selectedSubvolumeGroup } )}
		<DataTableLoading />
	{:then response}
		{@const subvolumes = response.subvolumes}
		<DataTable
			{selectedScope}
			{selectedFacility}
			{selectedVolume}
			{selectedSubvolumeGroup}
			{subvolumes}
		/>
	{:catch}
		<DataTableLoading />
	{/await}
</main>
