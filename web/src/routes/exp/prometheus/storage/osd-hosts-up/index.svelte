<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import OutNumber from './down-number.svelte';
	import Rate from './rate.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Card.Root class="col-span-1 h-full w-full border-none bg-muted/40 shadow-none">
	<Card.Header class="h-[150px]">
		<Card.Title>
			<h1 class="text-nowrap text-3xl">Up</h1>
			<Badge>Hosts</Badge>
		</Card.Title>
		<Card.Description></Card.Description>
	</Card.Header>
	<Card.Content class="h-[200px]">
		<Rate {client} {scope} />
	</Card.Content>
	<Card.Footer>
		<OutNumber {client} {scope} />
	</Card.Footer>
</Card.Root>
