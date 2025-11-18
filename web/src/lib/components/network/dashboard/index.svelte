<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';

	import AvailableIPs from './overview/available-ips.svelte';
	import DHCP from './overview/dhcp.svelte';
	import Discovery from './overview/discovery.svelte';
	import DNSServer from './overview/dns-server.svelte';
	import NetworkTraffic from './overview/network-traffic.svelte';
	import NetworkTrafficByTime from './overview/network-traffic-by-time.svelte';

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let isReloading = $state(true);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	onMount(async () => {
		try {
			environmentService
				.getPrometheus({})
				.then((response) => {
					prometheusDriver = new PrometheusDriver({
						endpoint: `${env.PUBLIC_API_URL}/prometheus`,
						baseURL: response.baseUrl
					});
				})
				.catch((error) => {
					console.error('Failed to initialize Prometheus driver:', error);
				});
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	});

	onDestroy(() => {
		isReloading = false;
	});
</script>

<main class="space-y-4 py-4">
	{#if prometheusDriver}
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
				<Tabs.Content value="overview">
					<div class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10">
						<div class="col-span-2">
							<Discovery bind:isReloading />
						</div>
						<div class="col-span-2">
							<DHCP bind:isReloading />
						</div>
						<div class="col-span-2 row-span-2">
							<AvailableIPs bind:isReloading />
						</div>
						<div class="col-span-4 row-span-2">
							<NetworkTrafficByTime {prometheusDriver} bind:isReloading />
						</div>
						<div class="col-span-2">
							<DNSServer bind:isReloading />
						</div>
						<div class="col-span-4 row-span-2">
							<NetworkTraffic {prometheusDriver} bind:isReloading />
						</div>
					</div>
				</Tabs.Content>
				<Tabs.Content value="analytics"></Tabs.Content>
			</Tabs.Root>
		</div>
	{/if}
</main>
