<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const CHART_TITLE = m.ram();

	// Queries
	const queries = $derived({
		description: `sum(node_memory_MemTotal_bytes{instance=~"${machine.fqdn}"})`,
		usage: `
            1 - (
                sum(node_memory_MemAvailable_bytes{instance=~"${machine.fqdn}"}) /
                sum(node_memory_MemTotal_bytes{instance=~"${machine.fqdn}"})
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
		const usagePercentage = usageValue != null ? usageValue * 100 : null;

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
{:catch}
	<ErrorLayout title={CHART_TITLE} />
{/await}
