<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const CHART_TITLE = 'CPU';
	const DESCRIPTION_UNIT = 'Cores';

	// Queries
	const queries = $derived({
		description: `count(count by (cpu, instance) (node_cpu_seconds_total{instance=~"${machine.fqdn}"}))`,
		usage: `
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}",mode!="idle"}[6m]))
			/
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}"}[6m]))
		`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [descriptionResponse, usageResponse] = await Promise.all([
			client.instantQuery(queries.description),
			client.instantQuery(queries.usage)
		]);

		// Parse responses
		const description = descriptionResponse.result[0]?.value?.value ?? null;
		const usageValue = usageResponse.result[0]?.value?.value ?? null;

		// Format data for display
		const formattedDescription = description ? `${description} ${DESCRIPTION_UNIT}` : undefined;
		const usagePercentage = usageValue != null ? usageValue * 100 : null;

		return {
			description: formattedDescription,
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
			{#if response.description}
				<Description description={response.description} />
			{:else}
				<Description />
			{/if}
		{/snippet}

		{#snippet content()}
			<Content data={response.usage} />
		{/snippet}
	</Layout>
{:catch error}
	Error
{/await}
