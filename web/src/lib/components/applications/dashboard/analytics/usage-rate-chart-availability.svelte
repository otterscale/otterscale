<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/standard.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.kubelet_availability();
	const CHART_DESCRIPTION = m.kubelet_api_server();

	// Queries
	const queries = $derived({
		usage: `apiserver_request:availability30d{juju_model_uuid=~"${scope.uuid}",verb="all"}`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [usageResponse] = await Promise.all([client.instantQuery(queries.usage)]);
		const usageValue = usageResponse.result[0]?.value?.value;
		const usagePercentage = usageValue != null ? usageValue * 100 : null;
		return {
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
			<Content data={response.usage} />
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} description={CHART_DESCRIPTION} />
{/await}
