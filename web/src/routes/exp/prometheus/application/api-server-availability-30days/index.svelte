<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Availability from './availability.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Template.Metric title="Availability">
	{#snippet description()}
		<p>API Server</p>
	{/snippet}
	{#snippet content()}
		<Availability {client} {scope} />
	{/snippet}
</Template.Metric>
