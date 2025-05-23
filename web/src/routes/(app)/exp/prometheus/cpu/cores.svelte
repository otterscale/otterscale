<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';

	let {
		client,
		juju_model_uuid,
		instance: instance
	}: { client: PrometheusDriver; juju_model_uuid: string; instance: string } = $props();
	const query = $derived(
		`
		count(
		count by (cpu) (node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${juju_model_uuid}"})
		)
		`
	);
</script>

{#await client.instantQuery(query)}
	<ComponentLoading />
{:then response}
	{@const cores = response.result[0].value.value}
	<p class="text-3xl">{cores} Cores</p>
{:catch error}
	Error
{/await}
