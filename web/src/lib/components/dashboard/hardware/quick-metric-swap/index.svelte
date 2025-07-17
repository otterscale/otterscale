<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Template from '../../utils/templates';
	import Memory from './memory.svelte';
	import FinalUsage from './usage.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="SWAP">
	{#snippet hint()}
		<p>Used Swap</p>
	{/snippet}
	{#snippet description()}
		<Memory {client} {scope} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {scope} />
	{/snippet}
</Template.Metric>
