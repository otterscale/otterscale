<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
	import { PrometheusDriver } from 'prometheus-query';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.cluster_health();

	// Query
	const query = $derived(
		`
		ceph_health_status{juju_model_uuid=~"${scope.uuid}"}
		`,
	);

	// Health status mappings
	const HEALTH_STATUS = {
		0: { label: 'HEALTHY', color: 'text-healthy' },
		1: { label: 'WARNING', color: 'text-warning' },
		2: { label: 'ERROR', color: 'text-error' },
	} as const;
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
				{@const value = result[0].value.value}
				{@const healthStatus = HEALTH_STATUS[value as keyof typeof HEALTH_STATUS]}
				<Content value={healthStatus?.label} textClass={healthStatus?.color} />
			{/if}
		{/snippet}
	</Layout>
{:catch error}
	<ErrorLayout title={CHART_TITLE} />
{/await}
