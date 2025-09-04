<script lang="ts">
	import { NetworkService, type Network } from '$lib/api/network/v1/network_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ArcChart, Text } from 'layerchart';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let { isReloading = $bindable() }: { isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');

	const networkClient = createClient(NetworkService, transport);

	const networks = writable<Network[]>([]);

	const targetSubnet = $derived($networks.find((network) => network?.vlan?.dhcpOn != null));

	const availableInternetProtocols = $derived([
		{
			key: 'available',
			value: Number(targetSubnet?.subnet?.statistics?.available ?? 0),
			color: 'var(--chart-2)',
		},
	]);
	const availableInternetProtocolsConfiguration = {} satisfies Chart.ChartConfig;

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
	<Card.Root class="h-full gap-2">
		<Card.Header class="items-center">
			<Card.Title>{m.available_ip_addresses()}</Card.Title>
			<Card.Description>
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger>
							{targetSubnet?.subnet?.statistics?.availablePercent}
						</Tooltip.Trigger>
						<Tooltip.Content>
							{targetSubnet?.subnet?.statistics?.available} / {targetSubnet?.subnet?.statistics?.total}
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			<Chart.Container
				config={availableInternetProtocolsConfiguration}
				class="mx-auto aspect-square max-h-[200px]"
			>
				<ArcChart
					label="key"
					value="value"
					outerRadius={88}
					innerRadius={66}
					trackOuterRadius={83}
					trackInnerRadius={72}
					padding={40}
					range={[90, -270]}
					maxValue={Number(targetSubnet?.subnet?.statistics?.total ?? 0)}
					series={availableInternetProtocols.map((ip) => ({
						key: ip.key,
						color: ip.color,
						data: [ip],
					}))}
					props={{
						arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
						tooltip: { context: { hideDelay: 350 } },
					}}
					tooltip={false}
				>
					{#snippet belowMarks()}
						<circle cx="0" cy="0" r="80" class="fill-background" />
					{/snippet}
					{#snippet aboveMarks()}
						<Text
							value={String(Number(targetSubnet?.subnet?.statistics?.available ?? 0))}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-4xl! font-bold"
							dy={3}
						/>
						<Text
							value={m.available()}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-muted-foreground!"
							dy={22}
						/>
					{/snippet}
				</ArcChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
