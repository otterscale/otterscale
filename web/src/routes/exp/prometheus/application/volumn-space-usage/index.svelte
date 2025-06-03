<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Usage from './usage.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Template.Metric title="Volumn Space">
	{#snippet content()}
		<Usage {client} {scope} />
	{/snippet}
</Template.Metric>
