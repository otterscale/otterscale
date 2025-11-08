<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.cpu();

	// Query
	const query = $derived(
		`
		avg(
			1
		-
			(
				rate(node_cpu_seconds_total{juju_model_uuid=~"${scope.uuid}",mode="idle"}[4m])
			or
				rate(node_cpu{juju_model_uuid=~"${scope.uuid}",mode="idle"}[4m])
			)
		)
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
				{@const value = result[0].value.value * 100}
				<Content value={value.toFixed(2)} unit="%" />
			{/if}
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} />
{/await}
