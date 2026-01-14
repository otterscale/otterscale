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
	async function fetchPods() {
        const allocateResponse = await prometheusDriver.instantQuery(
            `sum(kube_node_status_allocatable{resource="pods", juju_model="${scope}", container!=""})`
		);
        maxAllocatablePods = allocateResponse.result[0]?.value ?? undefined;

		const runningResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_status_phase{phase="Running", juju_model="${scope}", container!=""})`
		);
		runningPods = runningResponse.result[0]?.value ?? undefined;

		const pendingResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_status_phase{phase="Pending", juju_model="${scope}", container!=""})`
		);
		pendingPods = pendingResponse.result[0]?.value ?? undefined;
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
		icon="ph:package"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>Pods</Card.Title>
		<Card.Description class="flex h-6 items-center">Cluster Pods State</Card.Description>
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
		<Card.Content>
			<div class="flex gap-4">
                <div class="flex flex-col">
					<span class="text-3xl font-bold">{maxAllocatablePods?.value ?? 'N/A'}</span>
					<span class="text-1xl font-medium text-muted-foreground uppercase tracking-wider">Allocatable</span>
				</div>
				<div class="flex flex-col">
					<span class="text-3xl font-bold">{runningPods?.value ?? 'N/A'}</span>
					<span class="text-1xl font-medium text-muted-foreground uppercase tracking-wider">Running</span>
				</div>
				<div class="flex flex-col">
					<span class="text-3xl font-bold">{pendingPods?.value ?? 'N/A'}</span>
					<span class="text-1xl font-medium text-muted-foreground uppercase tracking-wider">Pending</span>
				</div>
			</div>
		</Card.Content>
	{/if}
</Card.Root>