<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.disk();

	// Query
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
		`,
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
				{@const value = result[0].value.value * 100}
				<Content value={value.toFixed(2)} unit={'%'} />
			{/if}
		{/snippet}
	</Layout>
{:catch error}
	<ErrorLayout title={CHART_TITLE} />
{/await}
