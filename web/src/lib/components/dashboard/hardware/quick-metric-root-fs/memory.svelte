<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Empty from '../../utils/empty';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
	const query = $derived(
		`
		node_filesystem_size_bytes{fstype!="rootfs",instance="${machine.fqdn}",mountpoint="/"}
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
		{@const memory = result[0].value.value}
		{@const capacity = formatCapacity(memory / 1024 / 1024)}
		<p class="text-3xl">{capacity.value} {capacity.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
