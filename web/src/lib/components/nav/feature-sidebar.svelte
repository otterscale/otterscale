<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { siteConfig } from '$lib/config/site';

	import Features from './features.svelte';
</script>

<Sheet.Root>
	<Sheet.Trigger asChild let:builder>
		<Button
			variant="outline"
			size="icon"
			class="shrink-0 border-none bg-header shadow-none md:hidden"
			builders={[builder]}
		>
			<Icon icon="ph:polygon" class="h-8 w-8" />
			<span class="sr-only">Toggle navigation menu</span>
		</Button>
	</Sheet.Trigger>
	<Sheet.Content side="left">
		<nav class="grid gap-4 whitespace-nowrap text-lg font-medium">
			<a href="##" class="flex items-center gap-4 text-lg font-semibold">
				<Icon icon="ph:polygon" class="h-8 w-8" />
				<span class="sr-only">{siteConfig.name}</span>
			</a>
			<Features />
		</nav>
	</Sheet.Content>
</Sheet.Root>
