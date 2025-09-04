<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';

	import Container from './overview/container.svelte';
	import ControlPlane from './overview/control-plane.svelte';
	import CPU from './overview/cpu.svelte';
	import Memory from './overview/memory.svelte';
	import NetworkTraffic from './overview/network-traffic.svelte';
	import Pod from './overview/pod.svelte';
	import ThroughtPut from './overview/throughput.svelte';
	import Worker from './overview/worker.svelte';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';
	import { activeScope } from '$lib/stores';

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
						baseURL: response.baseUrl,
					});
				})
				.catch((error) => {
					console.error('Failed to initialize Prometheus driver:', error);
				});
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	});
</script>

{#if prometheusDriver && $activeScope}
	<div class="mx-auto grid w-full gap-6">
		<div class="grid gap-1">
			<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.dashboard()}</h1>
			<p class="text-muted-foreground">
				{m.management_dashboard_description()}
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
				<div class="col-span-2">
					<ControlPlane scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-2">
					<Worker scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-3 row-span-2">
					<CPU {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-3 row-span-2">
					<Memory {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-2 col-start-1">
					<Pod {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-2">
					<Container {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-3 col-start-5">
					<NetworkTraffic {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
				<div class="col-span-3">
					<ThroughtPut {prometheusDriver} scope={$activeScope} bind:isReloading />
				</div>
			</Tabs.Content>
			<Tabs.Content value="analytics">
				<!-- {#if isMounted && prometheusDriver && $activeScope}
				<Dashboard client={prometheusDriver} scope={$activeScope} />
			{:else}
				<Loading />
			{/if} -->
			</Tabs.Content>
		</Tabs.Root>
	</div>
{/if}
