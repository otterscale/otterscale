<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import Content from '$lib/components/custom/chart/content/text/text.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/small-error.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.pools();

	// Query
	const query = $derived(
		`
		count(ceph_pool_metadata{juju_model_uuid=~"${scope.uuid}"})
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
				{@const value = result[0].value.value}
				<Content {value} />
			{/if}
		{/snippet}
	</Layout>
{:catch}
	<ErrorLayout title={CHART_TITLE} />
{/await}
