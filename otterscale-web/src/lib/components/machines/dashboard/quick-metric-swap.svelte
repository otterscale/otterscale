<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/quick.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { formatCapacity } from '$lib/formatter';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const CHART_TITLE = 'Swap';

	// Queries
	const queries = $derived({
		description: `node_memory_SwapTotal_bytes{instance=~"${machine.fqdn}"}`,
		usage: `
		(
			(
				node_memory_SwapTotal_bytes{instance=~"${machine.fqdn}"}
			-
				node_memory_SwapFree_bytes{instance=~"${machine.fqdn}"}
			)
		/
			(node_memory_SwapTotal_bytes{instance=~"${machine.fqdn}"})
		)
		`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [descriptionResponse, usageResponse] = await Promise.all([
			client.instantQuery(queries.description),
			client.instantQuery(queries.usage)
		]);

		const descriptionValue = descriptionResponse.result[0]?.value?.value;
		const usageValue = usageResponse.result[0]?.value?.value;

		const capacity = descriptionValue ? formatCapacity(descriptionValue) : null;
		const usagePercentage = usageValue ? usageValue * 100 : null;

		return {
			description: capacity ? `${capacity.value} ${capacity.unit}` : undefined,
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
