<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Busy from './busy.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="CPU">
	{#snippet hint()}
		<div class="flex items-center gap-2">
			<Badge variant="outline" class="w-fit">average cpu busy</Badge>
			<p>Average CPU busy across all hosts (OSD, RGW, MON etc) within the cluster</p>
		</div>
	{/snippet}
	{#snippet content()}
		<Busy {client} {scope} />
	{/snippet}
</Template.Metric>
