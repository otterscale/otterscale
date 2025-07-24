<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages.js';
	import { getIconFromUrl } from './icon';
	import { page } from '$app/state';

	type SubItem = {
		title: string;
		url: string;
	};

	type MainItem = {
		title: string;
		url: string;
		items?: SubItem[];
	};

	let { items }: { items: MainItem[] } = $props();

	function isItemActive(url: string): boolean {
		return page.url.pathname.startsWith(url);
	}

	function hasSubItems(item: MainItem): boolean {
		return Boolean(item.items?.length);
	}
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>{m.platform()}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each items as mainItem (mainItem.title)}
			<Collapsible.Root open={isItemActive(mainItem.url)}>
				{#snippet child({ props })}
					<Sidebar.MenuItem {...props}>
						<Sidebar.MenuButton tooltipContent={mainItem.title}>
							{#snippet child({ props })}
								<a href={mainItem.url} {...props}>
									<Icon icon={getIconFromUrl(mainItem.url)} />
									<span>{mainItem.title}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>

						{#if hasSubItems(mainItem)}
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
									{#each mainItem.items ?? [] as subItem (subItem.title)}
										<Sidebar.MenuSubItem>
											<Sidebar.MenuSubButton href={subItem.url}>
												<span>{subItem.title}</span>
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
