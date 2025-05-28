<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import Number from './number.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	let { client, scope: scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<Card.Root class="col-span-1 h-full w-full border-none bg-muted/40 shadow-none">
	<Card.Header class="h-[150px]">
		<Card.Title class="flex">
			<h1 class="text-nowrap text-3xl">Pools</h1>
		</Card.Title>
		<Card.Description></Card.Description>
	</Card.Header>
	<Card.Content class="h-[150px]">
		<Number {client} {scope} />
	</Card.Content>
	<Card.Footer></Card.Footer>
</Card.Root>
