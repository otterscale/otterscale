<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/arc/arc.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/quick.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// State
	let chartData = $state([{ value: 0 }]);
	let totalStorage: number | null = $state(null);
	let usagePercentage: number | null = $state(null);
	let loading = $state(true);
	let error = $state(false);

	// Prometheus queries
	const queries = $derived({
		total: `node_filesystem_size_bytes{fstype!="rootfs",instance=~"${machine.fqdn}",mountpoint="/"}`,
		usage: `
		1
		-
		(
			(
				node_filesystem_avail_bytes{fstype!="rootfs",instance=~"${machine.fqdn}",mountpoint="/"}
			)
			/
			node_filesystem_size_bytes{fstype!="rootfs",instance=~"${machine.fqdn}",mountpoint="/"}
		)
		`
	});

	async function fetchMetrics() {
		try {
			loading = true;
			error = false;

			const [totalResponse, usageResponse] = await Promise.all([
				client.instantQuery(queries.total),
				client.instantQuery(queries.usage)
			]);

			totalStorage = totalResponse.result[0]?.value?.value ?? null;
			const rawUsage = usageResponse.result[0]?.value?.value;
			usagePercentage = rawUsage ? rawUsage * 100 : null;

			chartData = usagePercentage !== null ? [{ value: usagePercentage }] : [{ value: 0 }];
		} catch (err) {
			error = true;
			console.error('Failed to fetch root filesystem metrics:', err);
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
			<Title title="Root FS" />
		{/snippet}

		{#snippet description()}
			{#if totalStorage === null}
				<Description />
			{:else}
				{@const capacity = formatCapacity(totalStorage)}
				<Description description="{capacity.value} {capacity.unit}" />
			{/if}
		{/snippet}

		{#snippet content()}
			<Content data={chartData} />
		{/snippet}
	</Layout>
{/if}
