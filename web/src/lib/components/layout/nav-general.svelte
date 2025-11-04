<script lang="ts">
	import Icon from '@iconify/svelte';

	import type { Route } from './routes';

	import { page } from '$app/state';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { pathDisabled, pathHidden, urlIcon } from '$lib/path';
	import { currentCeph, currentKubernetes } from '$lib/stores';

	let { title, routes }: { title: string; routes: Route[] } = $props();

	const isItemActive = (url: string): boolean => page.url.pathname.startsWith(url);
	const hasSubItems = (route: Route): boolean => Boolean(route.items?.length);
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>{title}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each routes as route (route.path.title)}
			{#if !pathHidden(page.params.scope, route.path.url)}
				<Collapsible.Root open={isItemActive(route.path.url)}>
					{#snippet child({ props })}
						<Sidebar.MenuItem {...props}>
							<Sidebar.MenuButton
								tooltipContent={route.path.title}
								aria-disabled={pathDisabled(
									$currentCeph?.name,
									$currentKubernetes?.name,
									page.params.scope,
									route.path.url,
								)}
							>
								{#snippet child({ props })}
									<a href={route.path.url} {...props}>
										<Icon icon={urlIcon(route.path.url)} />
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
											{#if !pathHidden(page.params.scope, subRoute.url)}
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton
														href={subRoute.url}
														aria-disabled={pathDisabled(
															$currentCeph?.name,
															$currentKubernetes?.name,
															page.params.scope,
															route.path.url,
														)}
													>
														<span>{subRoute.title}</span>
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											{/if}
										{/each}
									</Sidebar.MenuSub>
								</Collapsible.Content>
							{/if}
						</Sidebar.MenuItem>
					{/snippet}
				</Collapsible.Root>
			{/if}
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
