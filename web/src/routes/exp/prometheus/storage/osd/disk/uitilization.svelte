<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		avg(
			(
				(rate(node_disk_io_time_ms{juju_model_uuid=~"${scope.uuid}"}[4m]) / 1000)
			or
				(rate(node_disk_io_time_seconds_total{juju_model_uuid=~"${scope.uuid}"}[4m]))
			)
		* on (device) group_left (ceph_daemon)
			label_replace(
			ceph_disk_occupation_human{juju_model_uuid=~"${scope.uuid}"},
			"device",
			"$1",
			"device",
			"/dev/(.*)"
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
		{@const utilization = value * 100}
		<p class="text-5xl">{utilization.toFixed(2)}%</p>
	{/if}
{:catch error}
	Error
{/await}
