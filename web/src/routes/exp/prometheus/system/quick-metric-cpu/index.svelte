<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import FinalUsage from './usage.svelte';
	import Cores from './cores.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="CPU">
	{#snippet hint()}
		<p>Busy State of all CPU Cores</p>
	{/snippet}
	{#snippet description()}
		<Cores {client} {scope} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {scope} />
	{/snippet}
</Template.Metric>
