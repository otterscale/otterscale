<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import { formatCapacity } from '$lib/formatter';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import NoData from '../../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		count(ceph_pool_metadata{compression_mode!="none",juju_model_uuid=~"${scope.uuid}"})
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
			<Badge variant="outline" class="w-fit"># w/ compression</Badge>
		</span>
	{:else}
		{@const [result] = results}
		{@const number = result.value.value}
		<span class="flex flex-wrap items-end gap-2">
			<p class="text-xl">{number}</p>
			<Badge variant="outline" class="w-fit"># w/ compression</Badge>
		</span>
	{/if}
{:catch error}
	Error
{/await}
