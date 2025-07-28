<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Template from '../../utils/templates';
	import Memory from './memory.svelte';
	import FinalUsage from './usage.svelte';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
</script>

<Template.Metric title="SWAP">
	{#snippet hint()}
		<p>Used Swap</p>
	{/snippet}
	{#snippet description()}
		<Memory {client} {machine} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {machine} />
	{/snippet}
</Template.Metric>
