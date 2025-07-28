<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Template from '../../utils/templates';
	import Cores from './cores.svelte';
	import FinalUsage from './usage.svelte';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
</script>

<Template.Metric title="CPU">
	{#snippet hint()}
		<p>Busy State of all CPU Cores</p>
	{/snippet}
	{#snippet description()}
		<Cores {client} {machine} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {machine} />
	{/snippet}
</Template.Metric>
