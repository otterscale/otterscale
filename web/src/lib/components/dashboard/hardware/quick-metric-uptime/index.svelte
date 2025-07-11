<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Uptime from './uptime.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Text title="Uptime">
	{#snippet hint()}
		<p>Hardware Uptime</p>
	{/snippet}
	{#snippet content()}
		<Uptime {client} {scope} />
	{/snippet}
</Template.Text>
