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

	let nodeNotReady: SampleValue | undefined = $state(undefined);
	async function fetchNodeStatus() {
		const responseNotReadyNode = await prometheusDriver.instantQuery(
			`sum(kube_node_status_condition{condition="Ready" , status="false", juju_model="${scope}"})`
		);
		nodeNotReady = responseNotReadyNode.result[0]?.value ?? undefined;
	}

	async function fetch() {
		try {
			await fetchNodeStatus();
		} catch (error) {
			console.error('Failed to fetch cluster health:', error);
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
	<Card.Header>
		<Card.Title>{m.state()}</Card.Title>
		{#if nodeNotReady && Number(nodeNotReady.value) === 0}
			<Card.Description class="flex h-6 items-center">All Nodes Ready</Card.Description>
		{:else if nodeNotReady && Number(nodeNotReady.value) !== 0}
			<Card.Description class="flex h-6 items-center">
				{Number(nodeNotReady.value)} nodes unready
			</Card.Description>
		{/if}
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if nodeNotReady === undefined}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else if Number(nodeNotReady.value) === 0}
		<Card.Content class="text-3xl">Healthy</Card.Content>
	{:else}
		<Card.Content class="text-3xl">Unhealthy</Card.Content>
	{/if}
</Card.Root>
