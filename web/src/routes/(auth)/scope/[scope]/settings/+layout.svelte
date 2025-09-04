<script lang="ts">
	import { page } from '$app/state';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu';
	import { navigationMenuTriggerStyle } from '$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import Icon from '@iconify/svelte';

	let { children } = $props();

	let items = [
		{
			icon: 'ph:clock',
			title: m.ntp_server(),
			page: 'ntp-server',
		},
		{
			icon: 'ph:package',
			title: m.package_repository(),
			page: 'package-repository',
		},
		{
			icon: 'ph:disc',
			title: m.boot_image(),
			page: 'boot-image',
		},
		{
			icon: 'ph:cube',
			title: m.helm_repository(),
			page: 'helm-repository',
		},
		{
			icon: 'ph:tag-simple',
			title: m.machine_tag(),
			page: 'machine-tag',
		},
		{
			icon: 'ph:key',
			title: m.sso(),
			page: 'single-sign-on',
		},
		{
			icon: 'ph:test-tube',
			title: m.built_in_test(),
			page: 'built-in-test',
		},
		{
			icon: 'ph:wallet',
			title: m.subscription(),
			page: 'subscription',
		},
	];
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.settings()}</h1>
		<p class="text-muted-foreground">{m.settings_description()}</p>
	</div>

	<Separator />

	<div class="mx-auto grid w-full items-start gap-6 md:grid-cols-[180px_1fr] lg:grid-cols-[250px_1fr]">
		<NavigationMenu.Root viewport={false}>
			<NavigationMenu.List class="flex-col items-start">
				{#each items as item}
					{@const url = `${dynamicPaths.settings(page.params.scope).url}/${item.page}`}
					<NavigationMenu.Item>
						<NavigationMenu.Link>
							{#snippet child()}
								{#if page.url.pathname === url}
									<a href={url} class="{navigationMenuTriggerStyle()} bg-muted gap-2 font-semibold">
										<Icon icon={item.icon + '-bold'} class="size-4" />
										<span> {item.title} </span>
									</a>{:else}
									<a href={url} class="{navigationMenuTriggerStyle()} gap-2 font-normal">
										<Icon icon={item.icon} class="size-4" />
										<span> {item.title} </span>
									</a>
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
