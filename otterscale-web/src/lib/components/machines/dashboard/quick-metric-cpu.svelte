<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/quick.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const CHART_TITLE = 'CPU';

	// Queries
	const queries = $derived({
		cores: `count(count by (cpu) (node_cpu_seconds_total{instance=~"${machine.fqdn}"}))`,
		usage: `
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}",mode!="idle"}[6m]))
			/
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}"}[6m]))
		`
	});

	// Data fetching function
	async function fetchMetrics() {
		const [coresResponse, usageResponse] = await Promise.all([
			client.instantQuery(queries.cores),
			client.instantQuery(queries.usage)
		]);

		const cores = coresResponse.result[0]?.value?.value ?? null;
		const usage = usageResponse.result[0]?.value?.value
			? usageResponse.result[0].value.value * 100
			: null;

		return {
			cores,
			usage: usage !== null ? [{ value: usage }] : [{ value: 0 }],
			description: cores !== null ? `${cores} Cores` : undefined
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
