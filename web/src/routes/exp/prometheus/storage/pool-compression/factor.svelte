<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../../utils/empty.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		sum(ceph_pool_compress_under_bytes{juju_model_uuid=~"${scope.uuid}"} > 0)
		/
		sum(ceph_pool_compress_bytes_used{juju_model_uuid=~"${scope.uuid}"} > 0)
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
			<Badge variant="outline" class="w-fit">factor</Badge>
		</span>
	{:else}
		{@const [result] = results}
		{@const factor = result.value.value}
		<span class="flex flex-wrap items-end gap-2">
			<p class="text-3xl">{factor}</p>
			<Badge variant="outline" class="w-fit">factor</Badge>
		</span>
	{/if}
{:catch error}
	Error
{/await}
