<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import type { Facility } from '$lib/api/facility/v1/facility_pb';
	import { FacilityService } from '$lib/api/facility/v1/facility_pb';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let { scope, isReloading = $bindable() }: { scope: Scope; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);

	const facilities = writable<Facility[]>([]);
	const worker = $derived(
		$facilities.find((facility) => facility.name.includes('kubernetes-worker') && facility.units.length > 0),
	);
	const workerUnits = $derived(worker?.units ?? []);
	const activeWorkerUnits = $derived(workerUnits.filter((unit) => unit.workloadStatus?.state === 'active') ?? []);

	async function fetch() {
		facilityClient.listFacilities({ scopeUuid: scope.uuid }).then((response) => {
			facilities.set(response.facilities);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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

{#if isLoading}
	Loading
{:else}
	<Card.Root class="relative h-full gap-2 overflow-hidden">
		<Icon
			icon="ph:cube"
			class="text-primary/5 absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
		/>
		<Card.Header>
			<Card.Title>{m.workers()}</Card.Title>
			<Card.Description>{m.ready()}</Card.Description>
		</Card.Header>
		<Card.Content class="text-3xl">
			{activeWorkerUnits.length} / {workerUnits.length}
		</Card.Content>
	</Card.Root>
{/if}
