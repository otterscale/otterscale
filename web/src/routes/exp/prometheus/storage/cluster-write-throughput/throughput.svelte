<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';
	import { formatNetworkIO } from '$lib/formatter';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(irate(ceph_osd_op_w_in_bytes{juju_model_uuid=~"${scope.uuid}"}[5m]))
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
		{@const throughput = formatNetworkIO(value)}
		<p class="text-5xl">{throughput.value.toFixed(1)} {throughput.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
