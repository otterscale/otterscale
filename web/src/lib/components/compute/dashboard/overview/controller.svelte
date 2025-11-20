<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { Facility } from '$lib/api/facility/v1/facility_pb';
	import { FacilityService } from '$lib/api/facility/v1/facility_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let readyControllers: SampleValue = $state({} as SampleValue);
	async function fetchReadyControllers() {
		const response = await prometheusDriver.instantQuery(
			`count(kubevirt_virt_controller_ready_status)`
		);
		readyControllers = response.result[0]?.value?.value ?? {};
	}

	let healthControllers: SampleValue = $state({} as SampleValue);
	async function fetchHealthControllers() {
		const response = await prometheusDriver.instantQuery(
			`sum(kubevirt_virt_controller_ready_status)`
		);
		healthControllers = response.result[0]?.value?.value ?? {};
	}

	let isLoaded = $state(false);
	async function fetch() {
		try {
			await Promise.all([fetchReadyControllers(), fetchHealthControllers()]);
			isLoaded = true;
		} catch (error) {
			console.error('Failed to fetch cpu data:', error);
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
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoaded}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Icon
			icon="ph:compass"
			class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
		/>
		<Card.Header>
			<Card.Title>{m.controllers()}</Card.Title>
			<Card.Description>{m.ready()}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			{healthControllers} / {readyControllers}
		</Card.Content>
	</Card.Root>
{/if}
