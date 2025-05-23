<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { formatDuration } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';

	let {
		client,
		juju_model_uuid,
		instance
	}: { client: PrometheusDriver; juju_model_uuid: string; instance: string } = $props();
	const query = $derived(
		`
		node_time_seconds{instance="${instance}",juju_model_uuid=~"${juju_model_uuid}"}
		-
		node_boot_time_seconds{instance="${instance}",juju_model_uuid=~"${juju_model_uuid}"}
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const uptime = response.result[0].value.value}
	{@const duration = formatDuration(uptime)}
	<span class="flex items-end gap-2">
		<p class="text-6xl">{duration.value.toPrecision(2)}</p>
		<p class="text-4xl">{duration.unit}</p>
	</span>
{:catch error}
	Error
{/await}
