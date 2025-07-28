<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import Uptime from './uptime.svelte';
	import * as Template from '../../utils/templates';

	let { client, machine }: { client: PrometheusDriver; machine: Machine } = $props();
</script>

<Template.Text title="Uptime">
	{#snippet hint()}
		<p>Hardware Uptime</p>
	{/snippet}
	{#snippet content()}
		<Uptime {client} {machine} />
	{/snippet}
</Template.Text>
