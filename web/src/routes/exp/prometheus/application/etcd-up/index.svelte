<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Number from './number.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let {
		client,
		scope: scope,
	}: { client: PrometheusDriver; scope: Scope} = $props();
</script>

<Template.Text title="UP">
	{#snippet description()}
		<p>ETCD</p>
	{/snippet}
	{#snippet content()}
		<Number {client} {scope} />
	{/snippet}
</Template.Text>
