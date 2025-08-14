<script lang="ts" module>
	import { type Application_Chart } from '$lib/api/application/v1/application_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { fuzzLogosIcon } from './utils';
</script>

<script lang="ts">
	let { chart }: { chart: Application_Chart } = $props();
</script>

<Card.Root
	class={cn(
		chart.deprecated ? 'bg-muted' : 'hover:shadow-lg',
		'relative flex h-full flex-col justify-between gap-4 overflow-hidden transition-all'
	)}
>
	<Card.Header>
		{#if chart.verified}
			<span class="absolute top-6 right-6">
				<Icon icon="ph:star-fill" class="h-4 w-4 fill-yellow-400 text-yellow-400" />
			</span>
		{/if}
		<div class="flex items-center gap-2">
			<Avatar.Root class="h-12 w-12">
				<Avatar.Image src={chart.icon} />
				<Avatar.Fallback>
					<Icon icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')} class="size-12" />
				</Avatar.Fallback>
			</Avatar.Root>
			<span>
				<h3 class="font-semibold">{chart.name}</h3>
				<p class="text-muted-foreground flex items-center gap-1 text-sm">
					{chart.versions[0].chartVersion}
				</p>
			</span>
		</div>
	</Card.Header>
	<Card.Content class="text-muted-foreground p-8 text-sm">
		<p class="line-clamp-3 text-left">{chart.description}</p>
	</Card.Content>
	<Card.Footer>
		<Badge variant="outline" class="text-muted-foreground text-sm">
			{chart.license}
		</Badge>
	</Card.Footer>
</Card.Root>
