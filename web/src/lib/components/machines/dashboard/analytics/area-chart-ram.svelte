<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/area/area.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatTimeRange } from '$lib/components/custom/chart/units/formatter';
	import { fetchMultipleFlattenedRange } from '$lib/components/custom/prometheus';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	// Constants
	const STEP_SECONDS = 60; // 1 minute step
	const TIME_RANGE_HOURS = 1; // 1 hour of data

	// Chart configuration
	const CHART_TITLE = m.ram();
	const CHART_DESCRIPTION = m.memory_usage();

	// Time range calculation
	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	// Prometheus query for Memory usage
	const query = $derived({
		Total: `sum(node_memory_MemTotal_bytes{instance=~"${fqdn}"}) - sum(node_memory_MemAvailable_bytes{instance=~"${fqdn}"})`,
		Buffer: `sum(node_memory_Buffers_bytes{instance=~"${fqdn}"})`,
		Cache: `sum(node_memory_Cached_bytes{instance=~"${fqdn}"})`,
		Free: `sum(node_memory_MemFree_bytes{instance=~"${fqdn}"})`
	});
</script>

{#await fetchMultipleFlattenedRange(client, query, startTime, endTime, STEP_SECONDS)}
	<ComponentLoading />
{:then response}
	<Layout>
		{#snippet title()}
			<Title title={CHART_TITLE} />
		{/snippet}

		{#snippet description()}
			<Description description={CHART_DESCRIPTION} />
		{/snippet}

		{#snippet content()}
			<Content
				data={response}
				timeRange={formatTimeRange(TIME_RANGE_HOURS)}
				valueFormatter={formatCapacity}
			/>
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
