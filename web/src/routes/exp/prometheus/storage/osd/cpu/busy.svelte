<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		avg(
			1
		-
			(
				rate(node_cpu_seconds_total{juju_model_uuid=~"${scope.uuid}",mode="idle"}[4m])
			or
				rate(node_cpu{juju_model_uuid=~"${scope.uuid}",mode="idle"}[4m])
			)
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const results = response.result}
	{#if results.length === 0}
		<NoData />
	{:else}
		{@const [result] = results}
		{@const value = result.value.value}
		{@const busy = value * 100}
		<p class="text-5xl">{busy.toFixed(2)}%</p>
	{/if}
{:catch error}
	Error
{/await}
