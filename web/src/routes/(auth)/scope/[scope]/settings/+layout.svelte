<script lang="ts">
	import Icon from '@iconify/svelte';

	import { page } from '$app/state';
	import * as NavigationMenu from '$lib/components/ui/navigation-menu';
	import { navigationMenuTriggerStyle } from '$lib/components/ui/navigation-menu/navigation-menu-trigger.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { cn } from '$lib/utils';

	let { children } = $props();

	const globalItems = [
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

	const scopeBasedItems = [
		{
			icon: 'ph:hard-drives',
			title: m.data_volume(),
			page: 'data-volume',
		},
		{
			icon: 'ph:cpu',
			title: m.instance_type(),
			page: 'instance-type',
		},
		{
			icon: 'ph:cube',
			title: m.extensions(),
			page: 'extensions',
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
		<NavigationMenu.Root viewport={false} class="flex-col items-start gap-4">
			<NavigationMenu.List class="w-full flex-col items-start gap-1">
				{#each globalItems as item}
					{@const url = `${dynamicPaths.settings(page.params.scope).url}/${item.page}`}
					<NavigationMenu.Item>
						<NavigationMenu.Link>
							{#snippet child()}
								<a
									href={url}
									class={cn(
										navigationMenuTriggerStyle(),
										'h-fit',
										page.url.pathname === url
											? 'bg-muted gap-2 font-semibold'
											: 'gap-2 font-normal',
									)}
								>
									<Icon
										icon={page.url.pathname === url ? item.icon + '-bold' : item.icon}
										class="size-5"
									/>
									{item.title}
								</a>
							{/snippet}
						</NavigationMenu.Link>
					</NavigationMenu.Item>
				{/each}
			</NavigationMenu.List>
			<Separator class="w-full" />
			<NavigationMenu.List class="w-full flex-col items-start gap-2">
				{#each scopeBasedItems as item}
					{@const url = `${dynamicPaths.settings(page.params.scope).url}/${item.page}`}
					<NavigationMenu.Item>
						<NavigationMenu.Link>
							{#snippet child()}
								<a
									href={url}
									class={cn(
										navigationMenuTriggerStyle(),
										'h-fit',
										page.url.pathname === url
											? 'bg-muted gap-2 font-semibold'
											: 'gap-2 font-normal',
									)}
								>
									<Icon
										icon={page.url.pathname === url ? item.icon + '-bold' : item.icon}
										class="size-5"
									/>
									{item.title}
								</a>
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
