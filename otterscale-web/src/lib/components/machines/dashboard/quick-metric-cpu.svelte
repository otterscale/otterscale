<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/quick.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// State
	let chartData = $state([{ value: 0 }]);
	let cores: number | null = $state(null);
	let usage: number | null = $state(null);
	let loading = $state(true);
	let error = $state(false);

	// Queries
	const queries = $derived({
		cores: `count(count by (cpu) (node_cpu_seconds_total{instance=~"${machine.fqdn}"}))`,
		usage: `
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}",mode!="idle"}[6m]))
			/
			sum(irate(node_cpu_seconds_total{instance=~"${machine.fqdn}"}[6m]))
		`
	});

	// Data fetching
	async function fetchMetrics() {
		try {
			loading = true;
			error = false;

			const [coresResponse, usageResponse] = await Promise.all([
				client.instantQuery(queries.cores),
				client.instantQuery(queries.usage)
			]);

			cores = coresResponse.result[0]?.value?.value ?? null;
			usage = usageResponse.result[0]?.value?.value
				? usageResponse.result[0].value.value * 100
				: null;

			if (usage !== null) {
				chartData = [{ value: usage }];
			}
		} catch (err) {
			error = true;
			console.error('Failed to fetch CPU metrics:', err);
		} finally {
			loading = false;
		}
	}

	onMount(fetchMetrics);
</script>

{#if loading}
	<ComponentLoading />
{:else if error}
	Error
{:else}
	<Layout>
		{#snippet title()}
			<Title title="CPU" />
		{/snippet}

		{#snippet description()}
			{#if cores === null}
				<Description />
			{:else}
				<Description description="{cores} Cores" />
			{/if}
		{/snippet}

		{#snippet content()}
			<Content data={chartData} />
		{/snippet}
	</Layout>
{/if}
