<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Eligibility from './eligibility.svelte';
	import Factor from './factor.svelte';
	import Savings from './savings.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Card.Root class="col-span-1 h-full w-full border-none bg-muted/40 shadow-none">
	<Card.Header class="h-[150px]">
		<Card.Title class="flex">
			<h1 class="text-nowrap text-3xl">Compression</h1>
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
						<Badge variant="outline" class="w-fit">factor</Badge>
						<p>
							This factor describes the average ratio of data eligible to be compressed divided by
							the data actually stored. It does not account for data written that was ineligible for
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
				</HoverCard.Content>
			</HoverCard.Root>
		</Card.Title>
		<Card.Description>
			<Factor {client} {scope} />
		</Card.Description>
	</Card.Header>
	<Card.Content class="h-[150px]">
		<Savings {client} {scope} />
	</Card.Content>
	<Card.Footer>
		<Eligibility {client} {scope} />
	</Card.Footer>
</Card.Root>
