<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Empty from '../../utils/empty';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
	const query = $derived(
		`
		count(
		count by (cpu) (node_cpu_seconds_total{instance="${machine.fqdn}"})
		)
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
		{@const cores = result[0].value.value}
		<p class="text-3xl">{cores} Cores</p>
	{/if}
{:catch error}
	Error
{/await}
