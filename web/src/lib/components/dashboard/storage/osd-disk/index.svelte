<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Utilization from './uitilization.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Disk">
	{#snippet hint()}
		<div class="flex items-center gap-2">
			<Badge variant="outline" class="w-fit">average disk utilization</Badge>
			<p>Average Disk utilization for all OSD data devices (i.e. excludes journal/WAL)</p>
		</div>
	{/snippet}
	{#snippet content()}
		<Utilization {client} {scope} />
	{/snippet}
</Template.Metric>
