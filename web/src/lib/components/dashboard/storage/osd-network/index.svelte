<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Load from './load.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Network">
	{#snippet hint()}
		<div class="flex items-center gap-2">
			<Badge variant="outline" class="w-fit">network load</Badge>
			<p>Total send/receive network load across all hosts in the ceph cluster</p>
		</div>
	{/snippet}
	{#snippet content()}
		<Load {client} {scope} />
	{/snippet}
</Template.Metric>
