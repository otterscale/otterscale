<script lang="ts" module>
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import type { Network } from '$lib/api/network/v1/network_pb';
	import DeleteFabric from './delete-fabric.svelte';
	import UpdateFabric from './update-fabric.svelte';
	import UpdateVLAN from './update-vlan.svelte';
	import UpdateSubnet from './update-subnet.svelte';
</script>

<script lang="ts">
	let {
		row,
		networks = $bindable()
	}: {
		row: Row<Network>;
		networks: Writable<Network[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	{#if row.original.fabric}
		<Layout.ActionItem>
			<UpdateFabric fabric={row.original.fabric} bind:networks />
		</Layout.ActionItem>
		{#if row.original.vlan}
			<Layout.ActionItem>
				<UpdateVLAN fabric={row.original.fabric} vlan={row.original.vlan} bind:networks />
			</Layout.ActionItem>
		{/if}
		{#if row.original.subnet}
			<Layout.ActionItem>
				<UpdateSubnet subnet={row.original.subnet} bind:networks />
			</Layout.ActionItem>
		{/if}
		<Layout.ActionItem>
			<DeleteFabric fabric={row.original.fabric} bind:networks />
		</Layout.ActionItem>
	{/if}
</Layout.Actions>
