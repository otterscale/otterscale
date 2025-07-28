<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { formatDuration } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import * as Empty from '../../utils/empty';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
	const query = $derived(
		`
		node_time_seconds{instance="${machine.fqdn}"}
		-
		node_boot_time_seconds{instance="${machine.fqdn}"}
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const result = response.result}
	{#if result.length === 0}
		<Empty.Text />
	{:else}
		{@const uptime = result[0].value.value}
		{@const duration = formatDuration(uptime)}
		<p class="text-6xl">{duration.value.toPrecision(2)}</p>
		<p class="text-4xl">{duration.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
