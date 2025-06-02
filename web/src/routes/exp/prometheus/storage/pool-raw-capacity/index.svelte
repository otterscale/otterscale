<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Consumed from './consumed.svelte';
	import LogicalStored from './logical-stored.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Capacity">
	{#snippet hint()}
		<div class="flex flex-col gap-2">
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">total</Badge>
				<p>Total raw capacity available to the cluster</p>
			</div>
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">consumed</Badge>
				<p>
					Total raw capacity consumed by user data and associated overheads (metadata + redundancy)
				</p>
			</div>
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">logical stored</Badge>
				<p>Total of client data stored in the cluster</p>
			</div>
		</div>
	{/snippet}
	{#snippet description()}
		<Badge>Consumed</Badge>
	{/snippet}
	{#snippet content()}
		<Consumed {client} {scope} />
	{/snippet}
	{#snippet footer()}
		<LogicalStored {client} {scope} />
	{/snippet}
</Template.Metric>
