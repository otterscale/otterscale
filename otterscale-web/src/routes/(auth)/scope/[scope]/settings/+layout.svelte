<script lang="ts">
	import { page } from '$app/state';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu';
	import { navigationMenuTriggerStyle } from '$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';

	let { children } = $props();

	let items = [
		{
			title: 'NTP Server',
			page: 'ntp-server'
		},
		{
			title: 'Package Repository',
			page: 'package-repository'
		},
		{
			title: 'Boot Image',
			page: 'boot-image'
		},
		{
			title: 'Machine Tag',
			page: 'machine-tag'
		},
		{
			title: m.sso(),
			page: 'single-sign-on'
		}
	];
</script>

<div class="mx-auto grid w-full max-w-6xl gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.settings()}</h1>
		<p class="text-muted-foreground">{m.settings_description()}</p>
	</div>

	<Separator />

	<div
		class="mx-auto grid w-full items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]"
	>
		<NavigationMenu.Root viewport={false}>
			<NavigationMenu.List class="flex-col items-start">
				{#each items as item}
					{@const url = `${dynamicPaths.settings(page.params.scope).url}/${item.page}`}
					<NavigationMenu.Item>
						<NavigationMenu.Link active={true} class="w-full">
							{#snippet child()}
								{#if page.url.pathname === url}
									<a
										href={url}
										class="{navigationMenuTriggerStyle()} bg-muted font-semibold"
										data-active={true}
									>
										{item.title}
									</a>
								{:else}
									<a href={url} class={navigationMenuTriggerStyle()}>{item.title}</a>
								{/if}
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
