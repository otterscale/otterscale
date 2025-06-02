<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import IOPS from './iops.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="IOPS">
	{#snippet description()}
		<Badge>Read</Badge>
	{/snippet}
	{#snippet content()}
		<IOPS {client} {scope} />
	{/snippet}
</Template.Metric>
