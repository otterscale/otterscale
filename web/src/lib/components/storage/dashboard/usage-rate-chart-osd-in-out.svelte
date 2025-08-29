<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.osd();
	const CHART_DESCRIPTION = m.osd_in_ratio();
	const CHART_FOOTER = m.osd_out();

	// Queries
	const queries = $derived({
		in: `sum(ceph_osd_in{juju_model_uuid=~"${scope.uuid}"})`,
		total: `count(ceph_osd_metadata{juju_model_uuid=~"${scope.uuid}"})`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [inResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.in),
			client.instantQuery(queries.total)
		]);

		const inValue = inResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;
		const outValue = totalValue - inValue;
		const inUsageValue = inValue / totalValue;

		const inUsagePercentage = inUsageValue != null ? inUsageValue * 100 : null;
		return {
			inNumber: inValue,
			outNumber: outValue,
			totalNumber: totalValue,
			inUsage: inUsagePercentage !== null ? [{ value: inUsagePercentage }] : [{ value: NaN }]
		};
	}
</script>

{#await fetchMetrics()}
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
			<Content data={response.inUsage} subtitle={`${response.inNumber}/${response.totalNumber}`} />
		{/snippet}

		{#snippet footer()}
			{@const outNumber = Number(response.outNumber)}
			<Badge variant="outline">{outNumber} {CHART_FOOTER}</Badge>
		{/snippet}
	</Layout>
{:catch error}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
