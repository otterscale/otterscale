<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let maxAllocatablePods: SampleValue | undefined = $state(undefined);
	let runningPods: SampleValue | undefined = $state(undefined);
	let pendingPods: SampleValue | undefined = $state(undefined);
	let failedPods: SampleValue | undefined = $state(undefined);
	async function fetchPods() {
		const allocateResponse = await prometheusDriver.instantQuery(
			`sum(kube_node_status_allocatable{resource="pods", juju_model="${scope}"})`
		);
		maxAllocatablePods = allocateResponse.result[0]?.value ?? undefined;

		const runningResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_status_phase{phase="Running", juju_model="${scope}"})`
		);
		runningPods = runningResponse.result[0]?.value ?? undefined;

		const pendingResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_status_phase{phase="Pending", juju_model="${scope}"})`
		);
		pendingPods = pendingResponse.result[0]?.value ?? undefined;

		const failedResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_status_phase{phase="Failed", juju_model="${scope}"})`
		);
		failedPods = failedResponse.result[0]?.value ?? undefined;
	}

	async function fetch() {
		try {
			await fetchPods();
		} catch (error) {
			console.error('Failed to fetch pods:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:cube"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.pods()}</Card.Title>
		<Card.Description class="flex h-6 items-center">
			{m.cluster_dashboard_pod_description()}
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !maxAllocatablePods && !runningPods && !pendingPods}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="my-4 flex flex-col gap-4">
			<div>
				<p class="text-3xl font-bold">{maxAllocatablePods?.value ?? 'N/A'}</p>
				<p class="text-1xl font-medium tracking-wider uppercase">Allocatable</p>
			</div>
			<div class="grid grid-cols-3">
				<div class="text-chart-2">
					<p class="text-3xl font-bold">{runningPods?.value ?? 'N/A'}</p>
					<p class="text-1xl font-medium tracking-wider uppercase">Running</p>
				</div>
				<div class="text-muted-foreground">
					<p class="text-3xl font-bold">{pendingPods?.value ?? 'N/A'}</p>
					<p class="text-1xl font-medium tracking-wider uppercase">Pending</p>
				</div>
				<div class="text-chart-1">
					<p class="text-3xl font-bold">{failedPods?.value ?? 'N/A'}</p>
					<p class="text-1xl font-medium tracking-wider uppercase">Failed</p>
				</div>
			</div>
		</Card.Content>
	{/if}
</Card.Root>
