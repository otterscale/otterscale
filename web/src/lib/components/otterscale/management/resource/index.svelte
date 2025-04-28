<script lang="ts">
	import { page } from '$app/state';
	import ListNetworks from './network/index.svelte';
	import ListMachines from './machine/index.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import { type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';

	let {
		networks,
		machines
	}: {
		networks: Network[];
		machines: Machine[];
	} = $props();

	let activeTab = $state(page.url.hash ? page.url.hash : '#machine');
</script>

<main class="p-4">
	<Tabs.Root value={activeTab}>
		<Tabs.List class="grid w-fit grid-cols-2 rounded-sm">
			<Tabs.Trigger value="#machine" class="flex items-center gap-1">MACHINE</Tabs.Trigger>
			<Tabs.Trigger value="#network">NETWORK</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="#machine">
			<ListMachines {machines} />
		</Tabs.Content>
		<Tabs.Content value="#network">
			<ListNetworks {networks} />
		</Tabs.Content>
	</Tabs.Root>
</main>
