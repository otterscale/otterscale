<script lang="ts">
	import { NetworkService, type Network } from '$lib/api/network/v1/network_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let { isReloading = $bindable(), span }: { isReloading: boolean; span: string } = $props();

	const transport: Transport = getContext('transport');

	const networkClient = createClient(NetworkService, transport);

	const networks = writable<Network[]>([]);

	const targetSubnet = $derived($networks.find((network) => network?.vlan?.dhcpOn != null));

	function fetch() {
		networkClient
			.listNetworks({})
			.then((response) => {
				networks.set(response.networks);
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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
	Loading
{:else}
	<Card.Root class={cn('relative gap-2 overflow-hidden', span)}>
		<Card.Header>
			<Card.Title>{m.dhcp()}</Card.Title>
			<Card.Description>{targetSubnet?.subnet?.name}</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if !targetSubnet?.vlan?.dhcpOn}
				<p class="text-3xl text-green-600 dark:text-green-400">{m.on()}</p>
				<Icon
					icon="ph:check"
					class="text-primary/5 absolute top-4 -right-6 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
			{:else}
				<p class="text-3xl text-yellow-600 dark:text-yellow-400">{m.off()}</p>
				<Icon
					icon="ph:gps-slash"
					class="text-primary/5 absolute top-4 -right-6 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
			{/if}
		</Card.Content>
	</Card.Root>
{/if}
