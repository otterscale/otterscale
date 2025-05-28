<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';

	let {
		client,
		scope: scope,
		instance: instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
	const query = $derived(
		`
		count(
		count by (cpu) (node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"})
		)
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
		{@const cores = result[0].value.value}
		<p class="text-3xl">{cores} Cores</p>
	{/if}
{:catch error}
	Error
{/await}
