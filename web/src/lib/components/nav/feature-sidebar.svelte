<script lang="ts">
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { siteConfig } from '$lib/config/site';

	import Features from './features.svelte';
	import { cn } from '$lib/utils';
</script>

<Sheet.Root>
	<Sheet.Trigger
		class={cn(
			buttonVariants({ variant: 'outline', size: 'icon' }),
			'shrink-0 border-none bg-header shadow-none md:hidden [&_svg]:size-8'
		)}
	>
		<Icon icon="ph:polygon" />
		<span class="sr-only">Toggle navigation menu</span>
	</Sheet.Trigger>
	<Sheet.Content side="left">
		<nav class="grid gap-4 whitespace-nowrap text-lg font-medium">
			<a href="/" class="flex items-center gap-4 text-lg font-semibold [&_svg]:size-8">
				<Icon icon="ph:polygon" />
				<span class="sr-only">{siteConfig.name}</span>
			</a>
			<Features />
		</nav>
	</Sheet.Content>
</Sheet.Root>
