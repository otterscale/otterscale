<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../../utils/empty.svelte';
	import { formatCapacity } from '$lib/formatter';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		count(ceph_pool_metadata{juju_model_uuid=~"${scope.uuid}"})
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<NoData />
	{:else}
		{@const number = result[0].value.value}
		<p class="text-5xl">{number}</p>
	{/if}
{:catch error}
	Error
{/await}
