<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/area/area.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatTimeRange } from '$lib/components/custom/chart/units/formatter';
	import { fetchMultipleFlattenedRange } from '$lib/components/custom/prometheus';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const STEP_SECONDS = 60; // 1 minute step
	const TIME_RANGE_HOURS = 1; // 1 hour of data

	// Chart configuration
	const CHART_TITLE = 'IOPS';
	const CHART_DESCRIPTION = `${m.read()}/${m.write()}`;

	// Time range calculation
	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	// Prometheus query for Memory usage
	const query = $derived({
		Read: `sum(irate(ceph_osd_op_r{juju_model_uuid=~"${scope.uuid}"}[5m]))`,
		Write: `sum(irate(ceph_osd_op_w{juju_model_uuid=~"${scope.uuid}"}[5m]))`
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
