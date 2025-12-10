<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { type Release } from '$lib/api/application/v1/application_pb';
	import { type Chart } from '$lib/api/registry/v1/registry_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Item from '$lib/components/ui/item/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import { fuzzLogosIcon } from './utils';
</script>

<script lang="ts">
	let { chart, chartReleases }: { chart: Chart; chartReleases: Release[] | undefined } = $props();
</script>

<Card.Root class={cn(chart.deprecated ? 'bg-muted' : 'hover:shadow-lg', 'relative h-52 gap-0')}>
	<Card.Header>
		{#if chart.repositoryName.startsWith('otterscale/')}
			<span class="absolute top-6 right-6">
				<Icon icon="ph:star-fill" class="h-4 w-4 fill-yellow-400 text-yellow-400" />
			</span>
		{/if}
		<Item.Root class="p-0">
			<Avatar.Root class="m-2 size-8">
				<Avatar.Image src={chart.icon} class="object-contain" />
				<Avatar.Fallback>
					<Icon icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')} class="size-12" />
				</Avatar.Fallback>
			</Avatar.Root>
			<Item.Content class="text-left">
				<Item.Title>
					{chart.name}
				</Item.Title>
				<Item.Description>
					{chart.version}
				</Item.Description>
			</Item.Content>
		</Item.Root>
	</Card.Header>
	<Card.Content class="px-8 py-4 text-sm text-muted-foreground">
		<p class="line-clamp-3 text-left">{chart.description}</p>
	</Card.Content>
	<Card.Footer class="flex">
		<Badge variant="default" class={cn('ml-auto', chartReleases ? 'visible' : 'hidden')}>
			{m.installed()}
		</Badge>
	</Card.Footer>
</Card.Root>
