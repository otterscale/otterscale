<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { NetworkService, type Network } from '$lib/api/network/v1/network_pb';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { formatCapacity, formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleBand, scaleUtc } from 'd3-scale';
	import { ArcChart, BarChart, Highlight, Text, type ChartContextValue } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';
	import type { Writable } from 'svelte/store';
	import { writable } from 'svelte/store';
	import { fade } from 'svelte/transition';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.networking(page.params.scope)
	});

	let isReloading = $state(true);
	//
	const transport: Transport = getContext('transport');

	const networkClient = createClient(NetworkService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	const prometheusDriver: Writable<PrometheusDriver | null> = writable(null);

	const networks = writable<Network[]>([]);

	const targetSubnet = $derived($networks.find((network) => network?.vlan?.dhcpOn != null));

	let isDNSServersExpand = $state(false);

	const availableInternetProtocols = $derived([
		{
			key: 'available',
			value: Number(targetSubnet?.subnet?.statistics?.available ?? 0),
			color: 'var(--chart-2)'
		}
	]);
	const availableInternetProtocolsConfiguration = {} satisfies Chart.ChartConfig;

	let receives = $state([] as SampleValue[]);
	let transmits = $state([] as SampleValue[]);
	let latestReceive = $state({} as number);
	let latestTransmit = $state({} as number);
	let activeTraffic = $state<keyof typeof trafficsConfigurations>('receive');
	let trafficsContext = $state<ChartContextValue>();
	let trafficsByTimeContext = $state<ChartContextValue>();
	const traffics = $derived(
		receives.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmits[index]?.value ?? 0
		}))
	);
	const latestTraffics = $derived({
		receive: latestReceive,
		transmit: latestTransmit
	});
	const trafficsConfigurations = {
		views: { label: 'Traffic', color: '' },
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;
	const activeTrafficConfiguration = $derived([
		{
			key: activeTraffic,
			label: trafficsConfigurations[activeTraffic].label,
			color: trafficsConfigurations[activeTraffic].color
		}
	]);

	let receivesByTime = $state([] as SampleValue[]);
	let transmitsByTime = $state([] as SampleValue[]);
	const trafficsByTime = $derived(
		receivesByTime.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmitsByTime[index]?.value ?? 0
		}))
	);
	const trafficsByTimeConfiguration = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	function fetch() {
		environmentService.getPrometheus({}).then((response) => {
			prometheusDriver.set(
				new PrometheusDriver({
					endpoint: `${env.PUBLIC_API_URL}/prometheus`,
					baseURL: response.baseUrl
				})
			);

			if ($prometheusDriver && $activeScope) {
				$prometheusDriver
					.rangeQuery(
						`sum(irate(node_network_receive_bytes_total{juju_model_uuid="${$activeScope.uuid}"}[4m]))`,
						new Date().setMinutes(0, 0, 0) - 24 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60
					)
					.then((response) => {
						receives = response.result[0].values;
					});
				$prometheusDriver
					.rangeQuery(
						`sum(irate(node_network_transmit_bytes_total{juju_model_uuid="${$activeScope.uuid}"}[4m]))`,
						new Date().setMinutes(0, 0, 0) - 24 * 60 * 60 * 1000,
						new Date().setMinutes(0, 0, 0),
						2 * 60
					)
					.then((response) => {
						transmits = response.result[0].values;
					});
				$prometheusDriver
					.instantQuery(
						`sum(irate(node_network_receive_bytes_total{juju_model_uuid="${$activeScope.uuid}"}[4m]))`
					)
					.then((response) => {
						latestReceive = response.result[0].value.value;
					});
				$prometheusDriver
					.instantQuery(
						`sum(irate(node_network_transmit_bytes_total{juju_model_uuid="${$activeScope.uuid}"}[4m]))`
					)
					.then((response) => {
						latestTransmit = response.result[0].value.value;
					});
				$prometheusDriver
					.rangeQuery(
						`sum(increase(node_network_receive_bytes_total{juju_model_uuid="${$activeScope.uuid}"}[1h]))`,
						new Date().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
						new Date().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
						1 * 60 * 60
					)
					.then((response) => {
						receivesByTime = response.result[0]?.values;
					});
				$prometheusDriver
					.rangeQuery(
						`sum(increase(node_network_transmit_bytes_total{}[1h]))`,
						new Date().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
						new Date().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
						1 * 60 * 60
					)
					.then((response) => {
						transmitsByTime = response.result[0]?.values;
					});
			}
		});

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

	onMount(async () => {
		fetch();
		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.networking()}</h1>
		<p class="text-muted-foreground">
			{m.networking_dashboard_description()}
		</p>
	</div>
	<Tabs.Root value="overview">
		<div class="flex justify-between gap-2">
			<Tabs.List>
				<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
				<Tabs.Trigger value="analytics" disabled>{m.analytics()}</Tabs.Trigger>
			</Tabs.List>
			<Reloader bind:checked={isReloading} />
		</div>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10"
		>
			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
				<Card.Header>
					<Card.Title>{m.discovery()}</Card.Title>
					<Card.Description>{targetSubnet?.subnet?.name}</Card.Description>
				</Card.Header>
				<Card.Content>
					{#if !targetSubnet?.subnet?.activeDiscovery}
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

			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
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

			<Card.Root class="col-span-2 row-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>{m.available_ip_addresses()}</Card.Title>
					<Card.Description>
						{targetSubnet?.subnet?.statistics?.availablePercent}
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
								data: [ip]
							}))}
							props={{
								arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
								tooltip: { context: { hideDelay: 350 } }
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

			<Card.Root class="col-span-4 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>{m.total_upload_and_download()}</Card.Title>
					<Card.Description>
						<p class="lowercase">{m.over_each_time({ unit: m.hour() })}</p>
					</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={trafficsByTimeConfiguration} class="h-[200px] w-full">
						<BarChart
							bind:context={trafficsByTimeContext}
							data={trafficsByTime}
							xScale={scaleBand().padding(0.25)}
							x="time"
							axis="x"
							rule={false}
							series={[
								{
									key: 'receive',
									label: trafficsByTimeConfiguration.receive.label,
									color: trafficsByTimeConfiguration.receive.color,
									props: { rounded: 'bottom' }
								},
								{
									key: 'transmit',
									label: trafficsByTimeConfiguration.transmit.label,
									color: trafficsByTimeConfiguration.transmit.color
								}
							]}
							seriesLayout="stack"
							props={{
								bars: {
									stroke: 'none',
									initialY: trafficsByTimeContext?.height,
									initialHeight: 0,
									motion: {
										y: { type: 'tween', duration: 500, easing: cubicInOut },
										height: { type: 'tween', duration: 500, easing: cubicInOut }
									}
								},
								highlight: { area: false },
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
									ticks: 1
								}
							}}
							legend
						>
							{#snippet belowMarks()}
								<Highlight area={{ class: 'fill-muted' }} />
							{/snippet}
							{#snippet tooltip()}
								<Chart.Tooltip
									labelFormatter={(time: Date) => {
										return time.toLocaleDateString('en-US', {
											year: 'numeric',
											month: 'short',
											day: 'numeric',
											hour: 'numeric',
											minute: 'numeric'
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										{@const { value: io, unit } = formatCapacity(Number(value))}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{io} {unit}</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
						</BarChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="relative col-span-2 gap-2 overflow-hidden">
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

			<Card.Root class="col-span-4 row-span-2 gap-2">
				<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
					<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
						<Card.Title>{m.network_traffic()}</Card.Title>
					</div>
					<div class="flex">
						{#each ['receive', 'transmit'] as key (key)}
							{@const chart = key as keyof typeof trafficsConfigurations}
							{@const { value, unit } = formatIO(
								latestTraffics[key as keyof typeof latestTraffics]
							)}
							<button
								data-active={activeTraffic === chart}
								class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6"
								onclick={() => (activeTraffic = chart)}
							>
								<span class="text-muted-foreground text-xs">
									{trafficsConfigurations[chart].label}
								</span>
								<span class="flex items-end gap-1 text-lg leading-none font-bold sm:text-3xl">
									{value.toLocaleString()}
									<span class="text-xs">{unit}</span>
								</span>
							</button>
						{/each}
					</div>
				</Card.Header>
				<Card.Content class="px-6 pt-6">
					<Chart.Container config={trafficsConfigurations} class="aspect-auto h-[120px] w-full">
						<BarChart
							bind:context={trafficsContext}
							data={traffics}
							x="time"
							axis="x"
							series={activeTrafficConfiguration}
							props={{
								bars: {
									stroke: 'none',
									rounded: 'none',
									// use the height of the chart to animate the bars
									initialY: trafficsContext?.height,
									initialHeight: 0,
									motion: {
										y: { type: 'tween', duration: 500, easing: cubicInOut },
										height: { type: 'tween', duration: 500, easing: cubicInOut }
									}
								},
								highlight: { area: { fill: 'none' } },
								xAxis: {
									format: (v: Date) =>
										`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
									ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks()
								}
							}}
						>
							{#snippet belowMarks()}
								<Highlight area={{ class: 'fill-muted' }} />
							{/snippet}
							{#snippet tooltip()}
								<Chart.Tooltip
									nameKey="views"
									labelFormatter={(time: Date) => {
										return time.toLocaleDateString('en-US', {
											year: 'numeric',
											month: 'short',
											day: 'numeric',
											hour: 'numeric',
											minute: 'numeric'
										});
									}}
								>
									{#snippet formatter({ item, name, value })}
										{@const { value: io, unit } = formatIO(Number(value))}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{io} {unit}</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
						</BarChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="analytics"></Tabs.Content>
	</Tabs.Root>
</div>
