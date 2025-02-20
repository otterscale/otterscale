<script lang="ts">
	import Icon from '@iconify/svelte';

	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { siteConfig } from '$lib/config/site';
	import { cn } from '$lib/utils';
	import { features } from './features';
</script>

<Sheet.Root>
	<Sheet.Trigger asChild let:builder>
		<Button
			variant="outline"
			size="icon"
			class="shrink-0 border-none shadow-none md:hidden"
			builders={[builder]}
		>
			<Icon icon="line-md:chevron-double-right" class="h-8 w-8" />
			<span class="sr-only">Toggle navigation menu</span>
		</Button>
	</Sheet.Trigger>
	<Sheet.Content side="left">
		<nav class="grid gap-4 whitespace-nowrap text-lg font-medium">
			<a href="##" class="flex items-center gap-4 text-lg font-semibold">
				<Icon icon="line-md:chevron-double-right" class="h-8 w-8" />
				<span class="sr-only">{siteConfig.name}</span>
			</a>
			{#each features as feature}
				<a
					href={feature.path}
					class={cn(
						'transition-colors hover:text-foreground',
						page.url.pathname.startsWith(feature.path) ? 'text-foreground' : 'text-muted-foreground'
					)}
				>
					<span>{feature.name}</span>
				</a>
			{/each}
		</nav>
	</Sheet.Content>
</Sheet.Root>
