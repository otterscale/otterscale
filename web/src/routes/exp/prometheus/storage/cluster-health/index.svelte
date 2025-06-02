<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Status from './status.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Health">
	{#snippet content()}
		<Status {client} {scope} />
	{/snippet}
</Template.Metric>
