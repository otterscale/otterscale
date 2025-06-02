<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import FinalUsage from './usage.svelte';
	import Memory from './memory.svelte';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Template.Metric title="SWAP">
	{#snippet hint()}
		<p>Used Swap</p>
	{/snippet}
	{#snippet description()}
		<Memory {client} {scope} {instance} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {scope} {instance} />
	{/snippet}
</Template.Metric>
