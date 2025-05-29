<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Consumed from './used.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Card.Root class="col-span-1 h-full w-full border-none bg-muted/40 shadow-none">
	<Card.Header class="h-[150px]">
		<Card.Title>
			<div class="flex">
				<h1 class="overflow-hidden whitespace-nowrap text-3xl">Capacity</h1>
				<HoverCard.Root>
					<HoverCard.Trigger>
						<Button variant="ghost" size="icon" class="hover:bg-muted">
							<Icon icon="ph:info" />
						</Button>
					</HoverCard.Trigger>
					<HoverCard.Content
						class="flex w-fit max-w-[38w] flex-col gap-2 text-xs text-muted-foreground"
					>
						<div class="flex items-center gap-2">
							<Badge variant="outline" class="w-fit">total</Badge>
							<p>Total raw capacity available to the cluster</p>
						</div>
						<div class="flex items-center gap-2">
							<Badge variant="outline" class="w-fit">consumed</Badge>
							<p>
								Total raw capacity consumed by user data and associated overheads (metadata +
								redundancy)
							</p>
						</div>
						<div class="flex items-center gap-2">
							<Badge variant="outline" class="w-fit">logical stored</Badge>
							<p>Total of client data stored in the cluster</p>
						</div>
					</HoverCard.Content>
				</HoverCard.Root>
			</div>
			<Badge>Used</Badge>
		</Card.Title>
		<Card.Description></Card.Description>
	</Card.Header>
	<Card.Content class="h-[200px]">
		<Consumed {client} {scope} />
	</Card.Content>
	<Card.Footer></Card.Footer>
</Card.Root>
