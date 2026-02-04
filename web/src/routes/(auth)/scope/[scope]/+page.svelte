<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { Reloader } from '$lib/components/custom/reloader';
	import CpuUsage from '$lib/components/dashbaord/cluster/overview/cpu-usage.svelte';
	import Deployments from '$lib/components/dashbaord/cluster/overview/deployments.svelte';
	import GPUMemorUsage from '$lib/components/dashbaord/cluster/overview/gpu-memory-usage.svelte';
	import GPUUtilization from '$lib/components/dashbaord/cluster/overview/gpu-utilization.svelte';
	import Health from '$lib/components/dashbaord/cluster/overview/health.svelte';
	import Information from '$lib/components/dashbaord/cluster/overview/information.svelte';
	import MemoryUsage from '$lib/components/dashbaord/cluster/overview/memory-usage.svelte';
	import Nodes from '$lib/components/dashbaord/cluster/overview/nodes.svelte';
	import Pods from '$lib/components/dashbaord/cluster/overview/pods.svelte';
	import Uptime from '$lib/components/dashbaord/cluster/overview/uptime.svelte';
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

				<div class="flex justify-end">
					<Reloader bind:checked={isReloading} />
				</div>

				<div
					class="grid auto-rows-[minmax(140px,auto)] grid-cols-2 gap-4 pt-4 md:gap-6 lg:grid-cols-4"
				>
					<div class="col-span-1">
						<Health {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-1">
						<Uptime {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-2">
						<Information {prometheusDriver} {scope} bind:isReloading />
					</div>

					<div class="col-span-1">
						<Nodes {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-1">
						<Deployments {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-2">
						<Pods {prometheusDriver} {scope} bind:isReloading />
					</div>

					<div class="col-span-1">
						<CpuUsage {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-1">
						<MemoryUsage {prometheusDriver} {scope} bind:isReloading />
					</div>

					<div class="col-span-1">
						<GPUUtilization {prometheusDriver} {scope} bind:isReloading />
					</div>
					<div class="col-span-1">
						<GPUMemorUsage {prometheusDriver} {scope} bind:isReloading />
					</div>
				</div>
			</div>
		{/if}
	</main>
{/key}
