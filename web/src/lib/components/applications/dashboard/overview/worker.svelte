<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { Facility } from '$lib/api/facility/v1/facility_pb';
	import { FacilityService } from '$lib/api/facility/v1/facility_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let { scope, isReloading = $bindable() }: { scope: string; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);

	const facilities = writable<Facility[]>([]);
	const worker = $derived(
		$facilities.find(
			(facility) => facility.name.includes('kubernetes-worker') && facility.units.length > 0
		)
	);
	const workerUnits = $derived(worker?.units ?? []);
	const activeWorkerUnits = $derived(
		workerUnits.filter((unit) => unit.workloadStatus?.state === 'active') ?? []
	);

	async function fetch() {
		try {
			const response = await facilityClient.listFacilities({ scope: scope });
			facilities.set(response.facilities);
		} catch (error) {
			console.error('Failed to fetch facilities:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:cube"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>{m.workers()}</Card.Title>
		<Card.Description>{m.ready()}</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else}
		<Card.Content class="text-3xl">
			{activeWorkerUnits.length} / {workerUnits.length}
		</Card.Content>
	{/if}
</Card.Root>
