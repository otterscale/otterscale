<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import { formatCapacity } from '$lib/formatter';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Empty from '../../utils/empty';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		(
			sum(ceph_pool_compress_under_bytes{juju_model_uuid=~"${scope.uuid}"} > 0)
		/
			sum(
				ceph_pool_stored_raw{juju_model_uuid=~"${scope.uuid}"}
			and
				ceph_pool_compress_under_bytes{juju_model_uuid=~"${scope.uuid}"} > 0
			)
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const results = response.result}
	{#if results.length === 0}
		<span class="flex w-full flex-wrap items-center justify-center gap-2">
			<Empty.Text class="w-fit" />
			<Badge variant="outline" class="w-fit">eligibility</Badge>
		</span>
	{:else}
		{@const [result] = results}
		{@const value = result.value.value}
		{@const eligibility = formatCapacity(value / 1024 / 1024)}
		<span class="flex flex-wrap items-end gap-2">
			<p class="text-xl">{eligibility.value} {eligibility.unit}</p>
			<Badge variant="outline" class="w-fit">eligibility</Badge>
		</span>
	{/if}
{:catch error}
	Error
{/await}
