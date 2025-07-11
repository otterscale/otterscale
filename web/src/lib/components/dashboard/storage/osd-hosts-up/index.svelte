<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import OutNumber from './down-number.svelte';
	import Rate from './rate.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Template from '../../utils/templates';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Template.Metric title="Up">
	{#snippet description()}
		<p class="text-xl">Hosts</p>
	{/snippet}
	{#snippet content()}
		<Rate {client} {scope} />
	{/snippet}
	{#snippet footer()}
		<OutNumber {client} {scope} />
	{/snippet}
</Template.Metric>
