<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';
	import { formatNetworkIO } from '$lib/formatter';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(irate(ceph_osd_op_r{juju_model_uuid=~"${scope.uuid}"}[4m]))
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
		{@const iops = result.value.value}
		<p class="text-5xl">{iops.toFixed(1)} ops/s</p>
	{/if}
{:catch error}
	Error
{/await}
