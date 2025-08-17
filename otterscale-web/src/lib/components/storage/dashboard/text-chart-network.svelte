<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';
	import { formatIO } from '$lib/formatter';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.network();

	// Query
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
	<Layout>
		{#snippet title()}
			<Title title={CHART_TITLE} />
		{/snippet}

		{#snippet content()}
			{@const result = response.result}
			{#if result.length === 0}
				<Content />
			{:else}
				{@const value = result[0].value.value}
				{@const throughput = formatIO(value)}
				<Content value={throughput.value} unit={throughput.unit} />
			{/if}
		{/snippet}
	</Layout>
{:catch error}
	Error
{/await}
