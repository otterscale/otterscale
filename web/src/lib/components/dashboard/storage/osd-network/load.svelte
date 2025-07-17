<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../../utils/empty';
	import { formatNetworkIO } from '$lib/formatter';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(
			(
				rate(node_network_receive_bytes{device!="lo",juju_model_uuid=~"${scope.uuid}"}[4m])
				or
				rate(node_network_receive_bytes_total{device!="lo",juju_model_uuid=~"${scope.uuid}"}[4m])
			)
			unless on (device, instance)
			label_replace(
				(bonding_slaves{juju_model_uuid=~"${scope.uuid}"} > 0),
				"device",
				"$1",
				"master",
				"(.+)"
			)
		)
		+
		sum(
			(
				rate(node_network_transmit_bytes{device!="lo",juju_model_uuid=~"${scope.uuid}"}[4m])
				or
				rate(node_network_transmit_bytes_total{device!="lo",juju_model_uuid=~"${scope.uuid}"}[4m])
			)
			unless on (device, instance)
			label_replace(
				(bonding_slaves{juju_model_uuid=~"${scope.uuid}"} > 0),
				"device",
				"$1",
				"master",
				"(.+)"
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
		<Empty.Text />
	{:else}
		{@const [result] = results}
		{@const value = result.value.value}
		{@const load = formatNetworkIO(value)}
		<p class="text-5xl">{load.value} {load.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
