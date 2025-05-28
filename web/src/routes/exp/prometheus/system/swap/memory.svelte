<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { formatCapacity } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';
	import NoData from '../../utils/empty.svelte';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
	const query = $derived(
		`
		node_memory_SwapTotal_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<NoData />
	{:else}
		{@const memory = result[0].value.value}
		{@const capacity = formatCapacity(memory / 1024 / 1024)}
		<p class="text-3xl">{capacity.value} {capacity.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
