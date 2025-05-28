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
            (
                rate(
                node_disk_reads_completed{juju_model_uuid=~"${scope.uuid}"}[4m]
                )
            or
                rate(
                node_disk_reads_completed_total{juju_model_uuid=~"${scope.uuid}"}[4m]
                )
            )
        +
            (
                rate(
                node_disk_writes_completed{juju_model_uuid=~"${scope.uuid}"}[4m]
                )
            or
                rate(
                node_disk_writes_completed_total{juju_model_uuid=~"${scope.uuid}"}[4m]
                )
            )
        )
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
		{@const iops = result[0].value.value}
		<p class="text-5xl">{Math.round(iops)}</p>
	{/if}
{:catch error}
	Error
{/await}
