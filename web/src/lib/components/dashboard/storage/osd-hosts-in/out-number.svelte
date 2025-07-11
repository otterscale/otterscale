<script lang="ts">
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';
	import { onMount } from 'svelte';
	import { fetchInstance } from '../../utils';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();

	const outQuery = $derived(
		`
        sum(ceph_osd_in{juju_model_uuid=~"${scope.uuid}"} == bool 0)
		`
	);

	let outResponse: InstantVector | undefined | null = $state();

	let mounted = $state(false);
	onMount(async () => {
		try {
			outResponse = await fetchInstance(client, outQuery);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if !mounted}
	<ComponentLoading />
{:else if !outResponse}
	<Empty.Text />
{:else}
	{@const outNumber = Number(outResponse.value.value)}
	<Badge variant="outline">{outNumber} Out</Badge>
{/if}
