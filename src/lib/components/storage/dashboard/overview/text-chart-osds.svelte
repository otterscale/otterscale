<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.osds();
	const CHART_DESCRIPTION = 'In & Up';

	// Queries
	const queries = $derived({
		in: `sum(ceph_osd_in{juju_model="${scope}"})`,
		up: `sum(ceph_osd_up{juju_model="${scope}"})`,
		total: `count(ceph_osd_metadata{juju_model="${scope}"})`
	});

	// Auto Update
	let response = $state({} as { inNumber: number; upNumber: number; totalNumber: number });
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Data fetching function
	async function fetch() {
		const [inResponse, upResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.in),
			client.instantQuery(queries.up),
			client.instantQuery(queries.total)
		]);

		const inValue = inResponse.result[0]?.value?.value;
		const upValue = upResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;

		response = {
			inNumber: inValue,
			upNumber: upValue,
			totalNumber: totalValue
		};
	}

	// Effects
	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:disc"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header class="items-center">
		<Card.Title>{CHART_TITLE}</Card.Title>
		<Card.Description>{CHART_DESCRIPTION}</Card.Description>
	</Card.Header>
	{#if isLoading}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if response === undefined}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="text-3xl">{`${response.inNumber} / ${response.upNumber}`}</Card.Content>
	{/if}
</Card.Root>
