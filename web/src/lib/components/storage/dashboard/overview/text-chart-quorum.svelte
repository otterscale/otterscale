<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable(),
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.quorum_status();
	const CHART_DESCRIPTION = 'In & Up';

	// Queries
	const queries = $derived({
		in: `sum(ceph_mon_quorum_status{juju_model_uuid=~"${scope.uuid}"})`,
		total: `
		count(ceph_mon_quorum_status{juju_model_uuid=~"${scope.uuid}"})
		`,
	});

	// Auto Update
	let response = $state({} as { inNumber: number; totalNumber: number });
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Data fetching function
	async function fetch() {
		const [inResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.in),
			client.instantQuery(queries.total),
		]);

		const inValue = inResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;

		response = {
			inNumber: inValue,
			totalNumber: totalValue,
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

	onMount(() => {
		fetch();
		isLoading = false;
	});
</script>

{#if isLoading}
	<ComponentLoading />
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">{`${response.inNumber} / ${response.totalNumber}`}</Card.Content>
	</Card.Root>
{/if}
