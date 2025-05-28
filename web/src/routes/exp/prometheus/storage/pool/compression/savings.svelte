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
		sum(
			ceph_pool_compress_under_bytes{juju_model_uuid=~"${scope.uuid}"}
		-
			ceph_pool_compress_bytes_used{juju_model_uuid=~"${scope.uuid}"}
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<span class="flex w-full flex-wrap items-center justify-center gap-2">
			<NoData class="w-fit" />
			<Badge variant="outline" class="w-fit">savings</Badge>
		</span>
	{:else}
		{@const value = result[0].value.value}
		{@const savings = formatCapacity(value / 1024 / 1024)}
		<span class="flex flex-wrap items-end gap-2">
			<p class="text-5xl">{savings.value} {savings.unit}</p>
			<Badge variant="outline" class="w-fit">savings</Badge>
		</span>
	{/if}
{:catch error}
	Error
{/await}
