<script lang="ts">
	import { NetworkService, type Network } from '$lib/api/network/v1/network_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { fade } from 'svelte/transition';

	let { isReloading = $bindable(), span }: { isReloading: boolean; span: string } = $props();

	let isDNSServersExpand = $state(false);

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
			<Card.Title>{m.dns_server()}</Card.Title>
			<Card.Description>{targetSubnet?.subnet?.name}</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if targetSubnet?.subnet?.dnsServers}
				{#if targetSubnet?.subnet?.dnsServers.length === 1}
					<div class="flex items-center gap-1">
						<Icon icon="ph:share-network" />
						<p class="text-sm">{targetSubnet?.subnet?.dnsServers[0]}</p>
					</div>
				{:else if targetSubnet?.subnet?.dnsServers.length > 1}
					<div class="flex flex-col gap-1">
						{#each targetSubnet?.subnet?.dnsServers as dnsServer, index}
							{#if index === 0}
								<div class="flex items-center gap-2">
									<div class="flex items-center gap-1">
										<Icon icon="ph:share-network" />
										<p class="text-sm">{dnsServer}</p>
									</div>
									<Button
										variant="outline"
										class="h-5 p-2 text-xs transition-all duration-300"
										onmouseenter={() => {
											isDNSServersExpand = true;
										}}
										onmouseleave={() => {
											isDNSServersExpand = false;
										}}
									>
										+ {targetSubnet?.subnet?.dnsServers.length - 1}
									</Button>
								</div>
							{:else if isDNSServersExpand}
								<div
									class="flex translate-y-0 items-center gap-1 opacity-100 transition-all duration-300"
									in:fade={{ duration: 200 }}
									out:fade={{ duration: 200 }}
								>
									<Icon icon="ph:share-network" />
									<p class="text-sm">{dnsServer}</p>
								</div>
							{/if}
						{/each}
					</div>
				{/if}
				<Icon
					icon="ph:network"
					class="text-primary/5 absolute top-3 -right-3 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				/>
			{/if}
		</Card.Content>
	</Card.Root>
{/if}
