<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(
			kubelet_running_containers{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
		)
		or
		sum(
			kubelet_running_container_count{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
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
		{@const number = result.value.value}
		<p class="text-6xl">{number}</p>
	{/if}
{:catch error}
	Error
{/await}
