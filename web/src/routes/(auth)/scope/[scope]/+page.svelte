<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import Consumption from '$lib/components/dashbaord/cluster/overview/consumption.svelte';
	import CpuUsage from '$lib/components/dashbaord/cluster/overview/cpu.svelte';
	import Deployments from '$lib/components/dashbaord/cluster/overview/deployments.svelte';
	import GPUMemoryUsage from '$lib/components/dashbaord/cluster/overview/gpu-memory.svelte';
	import GPUUtilization from '$lib/components/dashbaord/cluster/overview/gpu-utilization.svelte';
	import Health from '$lib/components/dashbaord/cluster/overview/health.svelte';
	import Information from '$lib/components/dashbaord/cluster/overview/information.svelte';
	import MemoryUsage from '$lib/components/dashbaord/cluster/overview/memory.svelte';
	import Nodes from '$lib/components/dashbaord/cluster/overview/nodes.svelte';
	import Pods from '$lib/components/dashbaord/cluster/overview/pods.svelte';
	import Uptime from '$lib/components/dashbaord/cluster/overview/uptime.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import { m } from '$lib/paraglide/messages';

	const scope = $derived(page.params.scope!);

	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let isReloading = $state(true);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	onMount(async () => {
		try {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
			});
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	});

	onDestroy(() => {
		isReloading = false;
	});
</script>

{#key scope}
	<main class="space-y-4 py-4">
		{#if prometheusDriver}
			<div class="mx-auto grid w-full gap-6">
				<div class="grid gap-1">
					<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.k8s_overview_title()}</h1>
					<p class="text-muted-foreground">
						{m.k8s_overview_description()}
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
						<div
							class="grid auto-rows-[minmax(140px,auto)] grid-cols-2 gap-4 pt-4 md:gap-6 lg:grid-cols-1 2xl:grid-cols-4"
						>
							<div class="col-span-1 row-span-1">
								<Health {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-1">
								<Uptime {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-1">
								<Information {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-1">
								<Consumption {prometheusDriver} {scope} bind:isReloading />
							</div>

							<div class="col-span-1 row-span-1">
								<Nodes {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-1">
								<Deployments {prometheusDriver} {scope} bind:isReloading />
							</div>

							<div class="col-span-2 row-span-1">
								<Pods {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-2">
								<CpuUsage {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-2">
								<MemoryUsage {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-2">
								<GPUMemoryUsage {prometheusDriver} {scope} bind:isReloading />
							</div>
							<div class="col-span-1 row-span-2">
								<GPUUtilization {prometheusDriver} {scope} bind:isReloading />
							</div>
						</div>
					</Tabs.Content>
					<Tabs.Content value="analytics"></Tabs.Content>
				</Tabs.Root>
			</div>
		{/if}
	</main>
{/key}
