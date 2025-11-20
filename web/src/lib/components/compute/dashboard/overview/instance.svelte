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

	let instances: SampleValue = $state({} as SampleValue);
	async function fetchInstances() {
		const response = await prometheusDriver.instantQuery(`count(kubevirt_info)`);
		instances = response.result[0]?.value?.value ?? {};
	}

	let isLoaded = $state(false);
	async function fetch() {
		try {
			await fetchInstances();
			isLoaded = true;
		} catch (error) {
			console.error('Failed to fetch instance data:', error);
		}
	}

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	const reloadManager = new ReloadManager(fetch);

	onMount(async () => {
		await fetch();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoaded}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Icon
			icon="ph:shipping-container"
			class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
		/>
		<Card.Header>
			<Card.Title>{m.instances()}</Card.Title>
			<Card.Description>{m.running()}</Card.Description>
		</Card.Header>
		<Card.Content class="text-6xl">
			{instances}
		</Card.Content>
	</Card.Root>
{/if}
