<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';

	import Uptime from './uptime.svelte';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';

	let {
		client,
		scope: scope,
		instance
	}: { client: PrometheusDriver; scope: Scope; instance: string } = $props();
</script>

<Card.Root class="col-span-1 w-full border-none bg-muted/40 shadow-none">
	<Card.Header class="h-[150px]">
		<Card.Title class="flex">
			<h1 class="text-3xl">Uptime</h1>
			<HoverCard.Root>
				<HoverCard.Trigger>
					<Button variant="ghost" size="icon" class="hover:bg-muted">
						<Icon icon="ph:info" />
					</Button>
				</HoverCard.Trigger>
				<HoverCard.Content class="w-fit max-w-[38w] text-xs text-muted-foreground">
					System Uptime
				</HoverCard.Content>
			</HoverCard.Root>
		</Card.Title>
		<Card.Description></Card.Description>
	</Card.Header>
	<Card.Content class="h-[200px]">
		<Uptime {client} {scope} {instance} />
	</Card.Content>
	<Card.Footer></Card.Footer>
</Card.Root>
