<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import FinalUsage from './usage.svelte';
	import Memory from './memory.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Template.Metric title="Root FS">
	{#snippet hint()}
		<p>Used Root FS</p>
	{/snippet}
	{#snippet description()}
		<Memory {client} {scope} {instance} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {scope} {instance} />
	{/snippet}
</Template.Metric>
