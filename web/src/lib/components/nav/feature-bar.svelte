<script lang="ts">
	import Icon from '@iconify/svelte';

	import { page } from '$app/state';
	import { siteConfig } from '$lib/config/site.js';
	import { cn } from '$lib/utils.js';
	import { features } from './features';
</script>

<nav
	class="hidden flex-col gap-6 whitespace-nowrap text-lg font-medium md:flex md:flex-row md:items-center md:gap-5 md:text-sm lg:gap-6"
>
	<a href="/" class="flex items-center gap-2 text-lg font-semibold md:text-base">
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
