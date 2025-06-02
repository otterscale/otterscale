<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import PhysicalIOPS from './physical-iops.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="IOPS">
	{#snippet hint()}
		<div class="flex items-center gap-2">
			<Badge variant="outline" class="w-fit">physical iops</Badge>
			<p>IOPS Load at the device as reported by the OS on all OSD hosts</p>
		</div>
	{/snippet}
	{#snippet description()}
		<Badge>Physical</Badge>
	{/snippet}
	{#snippet content()}
		<PhysicalIOPS {client} {scope} />
	{/snippet}
</Template.Metric>
