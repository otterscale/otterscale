<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { siteConfig } from '$lib/config/site';

	import { cn } from '$lib/utils';
	import { featureTitle, type Feature } from '../features';
	import { i18n } from '$lib/i18n';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	export let features: Feature[] = [];
</script>

<Sheet.Root>
	<Sheet.Trigger
		class={cn(
			buttonVariants({ variant: 'outline', size: 'icon' }),
			'relative flex aspect-square size-8 shrink-0 items-center justify-center gap-4 overflow-hidden rounded-lg bg-sidebar-primary text-lg font-semibold text-sidebar-primary-foreground md:hidden [&_svg]:size-6'
		)}
	>
		<Icon icon="ph:graph" />
		<span class="sr-only">Toggle navigation menu</span>
	</Sheet.Trigger>
	<Sheet.Content side="left">
		<nav class="grid gap-4 whitespace-nowrap text-lg font-medium">
			<a
				href="/"
				class="relative flex aspect-square size-8 shrink-0 items-center justify-center gap-4 overflow-hidden rounded-lg bg-sidebar-primary text-lg font-semibold text-sidebar-primary-foreground [&_svg]:size-6"
			>
				<Icon icon="ph:graph" />
				<span class="sr-only">{siteConfig.name}</span>
			</a>

			{#each features as feature}
				<Button
					variant="ghost"
					class={cn(
						'px-2 py-0 text-lg transition-colors',
						!feature.enable ? 'disabled:pointer-events-auto disabled:cursor-not-allowed' : '',
						i18n.route(page.url.pathname).startsWith(feature.path)
							? 'text-foreground'
							: 'text-muted-foreground'
					)}
					disabled={!feature.enable}
					onclick={() => goto(i18n.resolveRoute(feature.path))}
				>
					<span>{featureTitle(feature.path)}</span>
				</Button>
			{/each}
		</nav>
	</Sheet.Content>
</Sheet.Root>
