<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { formatDuration } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';
	import NoData from '../utils/empty.svelte';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
	const query = $derived(
		`
		node_time_seconds{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		-
		node_boot_time_seconds{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
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
		{@const uptime = result[0].value.value}
		{@const duration = formatDuration(uptime)}
		<span class="flex items-end gap-2">
			<p class="text-6xl">{duration.value.toPrecision(2)}</p>
			<p class="text-4xl">{duration.unit}</p>
		</span>
	{/if}
{:catch error}
	Error
{/await}
