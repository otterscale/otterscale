<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatCapacity } from '$lib/formatter';
	import { PrometheusDriver } from 'prometheus-query';
	import NoData from '../../utils/empty.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
	const query = $derived(
		`
		node_filesystem_size_bytes{fstype!="rootfs",juju_model_uuid=~"${scope.uuid}",mountpoint="/"}
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
		{@const memory = result[0].value.value}
		{@const capacity = formatCapacity(memory / 1024 / 1024)}
		<p class="text-3xl">{capacity.value} {capacity.unit}</p>
	{/if}
{:catch error}
	Error
{/await}
