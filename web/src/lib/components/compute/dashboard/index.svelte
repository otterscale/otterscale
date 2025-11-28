<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';

	import ExtensionsAlert from './extensions-alert.svelte';
	import Controller from './overview/controller.svelte';
	import CPU from './overview/cpu.svelte';
	import Instance from './overview/instance.svelte';
	import Memory from './overview/memory.svelte';
	import NetworkTraffic from './overview/network-traffic.svelte';
	import ThroughtPut from './overview/throughput.svelte';
	import Pod from './overview/virtual-machine.svelte';
	import Worker from './overview/worker.svelte';

	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let isReloading = $state(false);
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
	<ExtensionsAlert scope={page.params.scope!} />
	{#if prometheusDriver}
		<div class="mx-auto grid w-full gap-6">
			<div class="grid gap-1">
				<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.dashboard()}</h1>
				<p class="text-muted-foreground">
					{m.compute_dashboard_description()}
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
					class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-12"
				>
					<div class="col-span-2">
						<Controller {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-2">
						<Worker {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-4 row-span-2">
						<CPU {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-4 row-span-2">
						<Memory {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-2 col-start-1">
						<Pod {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-2">
						<Instance {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-6">
						<NetworkTraffic {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-6">
						<ThroughtPut {prometheusDriver} {scope} bind:isReloading />
					</div>
				</Tabs.Content>
				<Tabs.Content value="analytics"></Tabs.Content>
			</Tabs.Root>
		</div>
	{/if}
</main>
