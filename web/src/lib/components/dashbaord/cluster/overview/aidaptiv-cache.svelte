<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let aidaptivCacheTotal: SampleValue | undefined = $state(undefined);
	let aidaptivCacheUsed: SampleValue | undefined = $state(undefined);
	async function fetchAidaptivCache() {
		const totalResponse = await prometheusDriver.instantQuery(
			`sum(kube_node_status_capacity{resource="phison.com/aidaptivcache", juju_model="${scope}", container!=""})`
		);
		aidaptivCacheTotal = totalResponse.result[0]?.value ?? undefined;

		const usedResponse = await prometheusDriver.instantQuery(
			`sum(kube_pod_container_resource_requests{resource="phison.com/aidaptivcache", juju_model="${scope}", container!=""})`
		);
		aidaptivCacheUsed = usedResponse.result[0]?.value ?? undefined;
	}

	async function fetch() {
		try {
			await fetchAidaptivCache();
		} catch (error) {
			console.error('Failed to fetch Aidaptiv Cache:', error);
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
		icon="ph:disk"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>Phison Aidaptiv Cache</Card.Title>
		<Card.Description class="flex h-6 items-center">Used / Total</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !aidaptivCacheTotal && !aidaptivCacheUsed}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content>
			<div class="text-2xl font-bold">
				{aidaptivCacheUsed?.value ?? 'N/A'}
				<span class="text-sm font-normal text-muted-foreground">
					/ {aidaptivCacheTotal?.value ?? 'N/A'}
				</span>
			</div>
			{#if aidaptivCacheUsed && aidaptivCacheTotal && Number(aidaptivCacheTotal.value) > 0}
				<Progress
					value={(Number(aidaptivCacheUsed.value) / Number(aidaptivCacheTotal.value)) * 100}
					class="mt-2 h-2"
				/>
			{/if}
		</Card.Content>
	{/if}
</Card.Root>
