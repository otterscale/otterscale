<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let consumption: string | undefined = $state(undefined);
	async function fetchConsumption() {
		const response = await prometheusDriver.instantQuery(
			`sum(DCGM_FI_DEV_POWER_USAGE{juju_model="${scope}"})`
		);
		console.log(response);
		consumption = response.result[0]?.value?.value ?? undefined;
	}

	const reloadManager = new ReloadManager(fetchConsumption);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	let isLoaded = $state(false);
	onMount(async () => {
		await fetchConsumption();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:info"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.power_usage()}</Card.Title>
		<Card.Description class="text-md flex h-6 items-center"
			>{m.cluster_dashboard_power_usage_description()}</Card.Description
		>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !consumption}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content>
			<p class="text-3xl">{Number(consumption).toFixed(2)} Watt</p>
		</Card.Content>
	{/if}
</Card.Root>
