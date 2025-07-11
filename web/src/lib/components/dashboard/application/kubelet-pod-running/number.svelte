<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(kubelet_running_pods{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"})
		or
		sum(
			kubelet_running_pod_count{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const results = response.result}
	{#if results.length === 0}
		<Empty.Text />
	{:else}
		{@const [result] = results}
		{@const number = result.value.value}
		<p class="text-6xl">{number}</p>
	{/if}
{:catch error}
	Error
{/await}
