<script lang="ts">
	import Icon from '@iconify/svelte';

	import { page } from '$app/state';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu';
	import { navigationMenuTriggerStyle } from '$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import { m } from '$lib/paraglide/messages';
	import { activeScope } from '$lib/stores';
	import { cn } from '$lib/utils';

	import { items } from './data';

	let { children } = $props();
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.settings()}</h1>
		<p class="text-muted-foreground">
			{m.settings_description({ scope: $activeScope.name })}
		</p>
	</div>

	<Separator />

	<div
		class="mx-auto grid w-full items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
	>
		<NavigationMenu.Root viewport={false} class="flex-col items-start gap-4">
			<NavigationMenu.List class="w-full flex-col items-start gap-1">
				{#each items as item}
					<NavigationMenu.Item>
						<NavigationMenu.Link>
							{#snippet child()}
								<!-- eslint-disable svelte/no-navigation-without-resolve -->
								<a
									href={item.url}
									class={cn(
										navigationMenuTriggerStyle(),
										'h-fit',
										page.url.pathname === item.url
											? 'gap-2 bg-muted font-semibold'
											: 'gap-2 font-normal'
									)}
								>
									<Icon
										icon={page.url.pathname === item.url ? item.icon + '-fill' : item.icon}
										class="size-6"
									/>
									<div class="flex flex-col">
										<p class="text-sm">{item.title}</p>
										<p class="text-xs text-muted-foreground">{item.type}</p>
									</div>
								</a>
								<!-- eslint-enable svelte/no-navigation-without-resolve -->
							{/snippet}
						</NavigationMenu.Link>
					</NavigationMenu.Item>
				{/each}
			</NavigationMenu.List>
		</NavigationMenu.Root>

		<div class="grid gap-6">
			{@render children()}
		</div>
	</div>
</div>
