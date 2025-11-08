<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/area/area.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatTimeRange } from '$lib/components/custom/chart/units/formatter';
	import { fetchMultipleFlattenedRange } from '$lib/components/custom/prometheus';
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const STEP_SECONDS = 60; // 1 minute step
	const TIME_RANGE_HOURS = 1; // 1 hour of data

	// Chart configuration
	const CHART_TITLE = m.disk();
	const CHART_DESCRIPTION = `${m.read()}/${m.write()}`;

	// Time range calculation
	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	// Prometheus query for Disk Read/Write
	const query = $derived({
		Read: `sum (rate(node_disk_read_bytes_total{instance=~"${machine.fqdn}", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+)"}[5m]))`,
		Write: `sum (rate(node_disk_written_bytes_total{instance=~"${machine.fqdn}", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+)"}[5m]))`
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
				valueFormatter={formatIO}
			/>
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
