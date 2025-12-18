<script lang="ts" module>
	import type { Network } from '$lib/api/network/v1/network_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
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

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>{m.actions()}</Actions.Label>
	<Actions.Separator />
	{#if network.fabric}
		<Actions.Item>
			<UpdateFabric fabric={network.fabric} {reloadManager} closeActions={close} />
		</Actions.Item>
		{#if network.vlan}
			<Actions.Item>
				<UpdateVLAN
					fabric={network.fabric}
					vlan={network.vlan}
					{reloadManager}
					closeActions={close}
				/>
			</Actions.Item>
		{/if}
		{#if network.subnet}
			<Actions.Item>
				<UpdateSubnet subnet={network.subnet} {reloadManager} closeActions={close} />
			</Actions.Item>
		{/if}
		<Actions.Item>
			<DeleteFabric fabric={network.fabric} {reloadManager} closeActions={close} />
		</Actions.Item>
	{/if}
</Actions.List>
