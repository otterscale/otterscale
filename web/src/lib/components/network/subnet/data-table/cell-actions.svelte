<script lang="ts" module>
	import type { Network } from '$lib/api/network/v1/network_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import DeleteFabric from './action-delete-fabric.svelte';
	import UpdateFabric from './action-update-fabric.svelte';
	import UpdateSubnet from './action-update-subnet.svelte';
	import UpdateVLAN from './action-update-vlan.svelte';
</script>

<script lang="ts">
	let {
		network,
		reloadManager
	}: {
		network: Network;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	{#if network.fabric}
		<Layout.ActionItem>
			<UpdateFabric fabric={network.fabric} {reloadManager} />
		</Layout.ActionItem>
		{#if network.vlan}
			<Layout.ActionItem>
				<UpdateVLAN fabric={network.fabric} vlan={network.vlan} {reloadManager} />
			</Layout.ActionItem>
		{/if}
		{#if network.subnet}
			<Layout.ActionItem>
				<UpdateSubnet subnet={network.subnet} {reloadManager} />
			</Layout.ActionItem>
		{/if}
		<Layout.ActionItem>
			<DeleteFabric fabric={network.fabric} {reloadManager} />
		</Layout.ActionItem>
	{/if}
</Layout.Actions>
