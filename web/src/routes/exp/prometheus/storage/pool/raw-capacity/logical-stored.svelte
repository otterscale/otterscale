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
		sum(ceph_pool_stored{juju_model_uuid=~"${scope.uuid}"})
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const results = response.result}
	{#if results.length === 0}
		<span class="flex w-full flex-wrap items-center justify-center gap-2">
			<NoData class="w-fit" />
			<Badge variant="outline" class="w-fit">logical stored</Badge>
		</span>
	{:else}
		{@const [result] = results}
		{@const value = result.value.value}
		{@const capacity = formatCapacity(value / 1024 / 1024)}
		<span class="flex flex-wrap items-end gap-2">
			<p class="text-xl">{capacity.value} {capacity.unit}</p>
			<Badge variant="outline" class="w-fit">logical stored</Badge>
		</span>
	{/if}
{:catch error}
	Error
{/await}
