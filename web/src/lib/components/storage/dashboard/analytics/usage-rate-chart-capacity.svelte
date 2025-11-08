<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.capacity();
	const CHART_DESCRIPTION = 'Cluster';

	// Queries
	const queries = $derived({
		used: `ceph_cluster_total_used_bytes{juju_model_uuid=~"${scope.uuid}"}`,
		total: `ceph_cluster_total_bytes{juju_model_uuid=~"${scope.uuid}"}`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [usedResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.used),
			client.instantQuery(queries.total)
		]);

		const usedValue = usedResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;

		const usedCapacity = usedValue ? formatCapacity(usedValue) : null;
		const totalCapacity = totalValue ? formatCapacity(totalValue) : null;
		const usageValue = usedValue / totalValue;
		const usagePercentage = usageValue != null ? usageValue * 100 : null;

		return {
			used: usedCapacity ? `${usedCapacity.value} ${usedCapacity.unit}` : undefined,
			total: totalCapacity ? `${totalCapacity.value} ${totalCapacity.unit}` : undefined,
			usage: usagePercentage !== null ? [{ value: usagePercentage }] : [{ value: NaN }]
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
			<Content data={response.usage} subtitle={`${response.used}/${response.total}`} />
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
