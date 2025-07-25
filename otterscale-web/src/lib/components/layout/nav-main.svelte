<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { getIconFromUrl, type Path } from '$lib/path';

	type Route = {
		path: Path;
		items: Path[];
	};

	let { routes }: { routes: Route[] } = $props();

	const isItemActive = (url: string): boolean => page.url.pathname.startsWith(url);
	const hasSubItems = (route: Route): boolean => Boolean(route.items?.length);
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>{m.platform()}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each routes as route (route.path.title)}
			<Collapsible.Root open={isItemActive(route.path.url)}>
				{#snippet child({ props })}
					<Sidebar.MenuItem {...props}>
						<Sidebar.MenuButton tooltipContent={route.path.title}>
							{#snippet child({ props })}
								<a href={route.path.url} {...props}>
									<Icon icon={getIconFromUrl(route.path.url)} />
									<span>{route.path.title}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>

						{#if hasSubItems(route)}
							<Collapsible.Trigger>
								{#snippet child({ props })}
									<Sidebar.MenuAction {...props} class="data-[state=open]:rotate-90">
										<Icon icon="ph:caret-right" />
										<span class="sr-only">Toggle</span>
									</Sidebar.MenuAction>
								{/snippet}
							</Collapsible.Trigger>

							<Collapsible.Content>
								<Sidebar.MenuSub>
									{#each route.items as subRoute (subRoute.title)}
										<Sidebar.MenuSubItem>
											<Sidebar.MenuSubButton href={subRoute.url}>
												<span>{subRoute.title}</span>
											</Sidebar.MenuSubButton>
										</Sidebar.MenuSubItem>
									{/each}
								</Sidebar.MenuSub>
							</Collapsible.Content>
						{/if}
					</Sidebar.MenuItem>
				{/snippet}
			</Collapsible.Root>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
