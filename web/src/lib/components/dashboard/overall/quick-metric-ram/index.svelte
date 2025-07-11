<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import FinalUsage from './usage.svelte';
	import Memory from './memory.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="RAM">
	{#snippet hint()}
		<p>Non Available RAM memory</p>
	{/snippet}
	{#snippet description()}
		<Memory {client} {scope} />
	{/snippet}
	{#snippet content()}
		<FinalUsage {client} {scope} />
	{/snippet}
</Template.Metric>
