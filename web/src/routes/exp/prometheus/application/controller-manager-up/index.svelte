<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Number from './number.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Template.Text title="UP">
	{#snippet description()}
		<p>Controller Manager</p>
	{/snippet}
	{#snippet content()}
		<Number {client} {scope} {instance} />
	{/snippet}
</Template.Text>
