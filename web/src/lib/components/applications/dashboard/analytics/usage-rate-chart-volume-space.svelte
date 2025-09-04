<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.kubelet_volume_space();

	// Queries
	const queries = $derived({
		usage: `
		max without (instance, node) (
			(
				topk(
					1,
					kubelet_volume_stats_capacity_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
				-
				topk(
					1,
					kubelet_volume_stats_available_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
			)
			/
			topk(
				1,
				kubelet_volume_stats_capacity_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
			)
		*
			100
		)
		`,
	});

	// Data fetching function
	async function fetchMetrics() {
		const [usageResponse] = await Promise.all([client.instantQuery(queries.usage)]);
		const usageValue = usageResponse.result[0]?.value?.value;
		const usagePercentage = usageValue != null ? usageValue * 100 : null;
		return {
			usage: usagePercentage !== null ? [{ value: usagePercentage }] : [{ value: NaN }],
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

		{#snippet content()}
			<Content data={response.usage} />
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} />
{/await}
