<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button';
	import { siteConfig } from '$lib/config/site';

	import { cn } from '$lib/utils';
	import { featureTitle, type Feature } from '../features';
	import { i18n } from '$lib/i18n';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	export let features: Feature[] = [];
</script>

<nav
	class="hidden flex-col gap-6 whitespace-nowrap text-lg font-medium md:flex md:flex-row md:items-center md:gap-2 md:text-sm lg:gap-4"
>
	<a
		href="/"
		class="relative flex aspect-square size-8 shrink-0 items-center justify-center gap-4 overflow-hidden rounded-lg bg-sidebar-primary text-lg font-semibold text-sidebar-primary-foreground [&_svg]:size-6"
	>
		<Icon icon="ph:polygon" />
		<span class="sr-only">{siteConfig.name}</span>
	</a>
	{#each features as feature}
		<Button
			variant="ghost"
			class={cn(
				'px-2 py-0 transition-colors',
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
