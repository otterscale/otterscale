<script lang="ts">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatDuration } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();

	// Constants
	const CHART_TITLE = m.uptime();

	// Query
	const query = $derived(
		`
		node_time_seconds{instance=~"${machine.fqdn}"}
		-
		node_boot_time_seconds{instance=~"${machine.fqdn}"}
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	<Layout>
		{#snippet title()}
			<Title title={CHART_TITLE} />
		{/snippet}

		{#snippet content()}
			{@const result = response.result}
			{#if result.length === 0}
				<Content />
			{:else}
				{@const uptime = result[0].value.value}
				{@const duration = formatDuration(uptime)}
				<Content value={duration.value.toFixed(1)} unit={duration.unit} />
			{/if}
		{/snippet}
	</Layout>
{:catch error}
	Error
{/await}
