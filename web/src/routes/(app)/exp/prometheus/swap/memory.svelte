<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { formatCapacity } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';

	let {
		client,
		juju_model_uuid,
		instance
	}: { client: PrometheusDriver; juju_model_uuid: string; instance: string } = $props();
	const query = $derived(
		`
		node_memory_SwapTotal_bytes{instance="${instance}",juju_model_uuid=~"${juju_model_uuid}"}
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const memory = response.result[0].value.value}
	{@const capacity = formatCapacity(memory / 1024 / 1024)}
	<p class="text-3xl">{capacity.value} {capacity.unit}</p>
{:catch error}
	Error
{/await}
