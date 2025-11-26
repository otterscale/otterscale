<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let readyControllers: SampleValue | undefined = $state(undefined);
	async function fetchReadyControllers() {
		const response = await prometheusDriver.instantQuery(
			`count(kubevirt_virt_controller_ready_status{juju_model="${scope}"})`
		);
		readyControllers = response.result[0]?.value ?? undefined;
	}

	let healthControllers: SampleValue | undefined = $state(undefined);
	async function fetchHealthControllers() {
		const response = await prometheusDriver.instantQuery(
			`sum(kubevirt_virt_controller_ready_status{juju_model="${scope}"})`
		);
		healthControllers = response.result[0]?.value ?? undefined;
	}

	let isLoaded = $state(false);
	async function fetch() {
		try {
			await Promise.all([fetchReadyControllers(), fetchHealthControllers()]);
		} catch (error) {
			console.error('Failed to fetch controller data:', error);
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
		icon="ph:compass"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.controllers()}</Card.Title>
		<Card.Description>{m.ready()}</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !healthControllers || !readyControllers}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="text-3xl">
			{healthControllers.value} / {readyControllers.value}
		</Card.Content>
	{/if}
</Card.Root>
