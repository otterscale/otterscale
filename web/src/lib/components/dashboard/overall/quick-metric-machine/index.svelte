<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Instance from './instance.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Text title="Instance">
	{#snippet hint()}
		<p>Amount of Machines</p>
	{/snippet}
	{#snippet content()}
		<Instance {client} {scope} />
	{/snippet}
</Template.Text>
