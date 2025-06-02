<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import Eligibility from './eligibility.svelte';
	import Factor from './factor.svelte';
	import Savings from './savings.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Compression">
	{#snippet hint()}
		<div class="flex flex-col gap-2">
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">factor</Badge>
				<p>
					This factor describes the average ratio of data eligible to be compressed divided by the
					data actually stored. It does not account for data written that was ineligible for
					compression (too small, or compression yield too low)
				</p>
			</div>
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">savings</Badge>
				<p>
					A compression saving is determined as the data eligible to be compressed minus the
					capacity used to store the data after compression
				</p>
			</div>
			<div class="flex items-center gap-2">
				<Badge variant="outline" class="w-fit">eligibility</Badge>
				<p>
					Indicates how suitable the data is within the pools that are/have been enabled for
					compression - averaged across all pools holding compressed data
				</p>
			</div>
		</div>
	{/snippet}
	{#snippet description()}
		<Factor {client} {scope} />
	{/snippet}
	{#snippet content()}
		<Savings {client} {scope} />
	{/snippet}
	{#snippet footer()}
		<Eligibility {client} {scope} />
	{/snippet}
</Template.Metric>
