<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { SampleValue } from 'prometheus-query';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let virtualMachines: SampleValue = $state({} as SampleValue);
	async function fetchVirtualMachines() {
		const response = await prometheusDriver.instantQuery(
			`count(kubevirt_info{juju_model="${scope}"})`
		);
		virtualMachines = response.result[0]?.value ?? {};
	}

	let isLoaded = $state(false);
	async function fetchData() {
		try {
			await fetchVirtualMachines();
			isLoaded = true;
		} catch (error) {
			console.error('Failed to fetch virtual machine data:', error);
		}
	}

	const reloadManager = new ReloadManager(fetchData);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	onMount(async () => {
		await fetchData();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full gap-2 overflow-hidden">
	<Icon
		icon="ph:squares-four"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.virtual_machines()}</Card.Title>
		<Card.Description>{m.starting()}</Card.Description>
	</Card.Header>
	<Card.Content class="h-full">
		{#if !isLoaded}
			<div class="flex h-full w-full items-center justify-center">
				<Icon icon="svg-spinners:3-dots-bounce" class="size-8" />
			</div>
		{:else if !virtualMachines.value}
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{:else}
			<p class="text-6xl">{virtualMachines.value}</p>
		{/if}
	</Card.Content>
</Card.Root>
