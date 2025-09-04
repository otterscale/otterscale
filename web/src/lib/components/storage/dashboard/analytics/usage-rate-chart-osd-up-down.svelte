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
	const CHART_DESCRIPTION = m.osd_up_ratio();
	const CHART_FOOTER = m.osd_down();
	// Queries
	const queries = $derived({
		up: `sum(ceph_osd_up{juju_model_uuid=~"${scope.uuid}"})`,
		total: `count(ceph_osd_metadata{juju_model_uuid=~"${scope.uuid}"})`,
	});

	// Data fetching function
	async function fetchMetrics() {
		const [upResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.up),
			client.instantQuery(queries.total),
		]);

		const upValue = upResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;
		const downValue = totalValue - upValue;
		const upUsageValue = upValue / totalValue;

		const upUsagePercentage = upUsageValue != null ? upUsageValue * 100 : null;
		return {
			upNumber: upValue,
			downNumber: downValue,
			totalNumber: totalValue,
			upUsage: upUsagePercentage !== null ? [{ value: upUsagePercentage }] : [{ value: NaN }],
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
			<Content data={response.upUsage} subtitle={`${response.upNumber}/${response.totalNumber}`} />
		{/snippet}

		{#snippet footer()}
			{@const outNumber = Number(response.downNumber)}
			<Badge variant="outline">{outNumber} {CHART_FOOTER}</Badge>
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
