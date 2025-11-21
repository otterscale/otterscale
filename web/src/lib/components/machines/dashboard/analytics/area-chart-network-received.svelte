<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/area/area.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatTimeRange } from '$lib/components/custom/chart/units/formatter';
	import { fetchFlattenedRange } from '$lib/components/custom/prometheus';
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	// Constants
	const STEP_SECONDS = 60; // 1 minute step
	const TIME_RANGE_HOURS = 1; // 1 hour of data

	// Chart configuration
	const CHART_TITLE = m.networking();
	const CHART_DESCRIPTION = m.received();

	// Time range calculation
	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	// Prometheus query for CPU load average
	const query = $derived(
		`sum by (device) (rate(node_network_receive_bytes_total{instance=~"${fqdn}", device!="lo"}[5m]))`
	);
</script>

{#await fetchFlattenedRange(client, query, startTime, endTime, STEP_SECONDS)}
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
				valueFormatter={formatIO}
			/>
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
