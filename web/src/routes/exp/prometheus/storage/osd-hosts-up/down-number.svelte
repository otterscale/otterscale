<script lang="ts">
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';
	import { onMount } from 'svelte';
	import { fetchInstance } from '../../utils';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();

	const downQuery = $derived(
		`
        sum(ceph_osd_up{juju_model_uuid=~"${scope.uuid}"} == bool 0)
		`
	);

	let downResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			downResponse = await fetchInstance(client, downQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !downResponse}
	<NoData />
{:else}
	{@const downNumber = Number(downResponse.value.value)}
	<Badge variant="outline">{downNumber} Down</Badge>
{/if}
