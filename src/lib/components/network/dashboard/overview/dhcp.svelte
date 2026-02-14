<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { type Network, NetworkService } from '$lib/api/network/v1/network_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let { isReloading = $bindable() }: { isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');

	const networkClient = createClient(NetworkService, transport);

	const networks = writable<Network[]>([]);

	const targetSubnet = $derived($networks.find((network) => network?.vlan?.dhcpOn));

	async function fetch() {
		try {
			const response = await networkClient.listNetworks({});
			networks.set(response.networks);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if isLoading}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Card.Header class="h-[42px]">
			<Card.Title>{m.dhcp()}</Card.Title>
		</Card.Header>
		<Card.Content>
			<div class="flex h-9 w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-6" />
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Card.Header>
			<Card.Title>{m.dhcp()}</Card.Title>
			<Card.Description>{targetSubnet?.subnet?.name}</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if targetSubnet?.vlan?.dhcpOn}
				<p class="text-3xl text-green-600 dark:text-green-400">{m.on()}</p>
				<Icon
					icon="ph:check"
					class="absolute top-4 -right-6 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
				/>
			{:else}
				<p class="text-3xl text-yellow-600 dark:text-yellow-400">{m.off()}</p>
				<Icon
					icon="ph:gps-slash"
					class="absolute top-4 -right-6 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
				/>
			{/if}
		</Card.Content>
	</Card.Root>
{/if}
