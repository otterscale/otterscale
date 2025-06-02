<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Utilization from './uitilization.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="RAM">
	{#snippet hint()}
		<div class="flex items-center gap-2">
			<Badge variant="outline" class="w-fit">average ram utilization</Badge>
			<p>Average Memory Usage across all hosts in the cluster (excludes buffer/cache usage)</p>
		</div>
	{/snippet}

	{#snippet content()}
		<Utilization {client} {scope} />
	{/snippet}
</Template.Metric>
