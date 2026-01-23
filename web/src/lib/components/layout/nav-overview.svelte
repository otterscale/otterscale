<script lang="ts">
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import type { Component } from 'svelte';

	import { goto } from '$app/navigation';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/context.svelte.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { m } from '$lib/paraglide/messages';
	let {
		items
	}: {
		items: {
			name: string;
			url: string;
			icon: Component;
			edit: boolean;
		}[];
	} = $props();
	const sidebar = useSidebar();
	const increment = 5;
	let visibleCount = $state(3);

	const visibleItems = $derived(items.slice(0, visibleCount));
	const hasMoreItems = $derived(items.length > visibleCount);

	function showMoreItems(): void {
		visibleCount += increment;
	}
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>{m.overview()}</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each visibleItems as item (item.name)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<!-- eslint-disable svelte/no-navigation-without-resolve -->
						<a href={item.url} {...props}>
							<item.icon />
							<span>{item.name}</span>
						</a>
						<!-- eslint-enable svelte/no-navigation-without-resolve -->
					{/snippet}
				</Sidebar.MenuButton>
				{#if item.edit}
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Sidebar.MenuAction showOnHover {...props}>
									<EllipsisIcon />
									<span class="sr-only">More</span>
								</Sidebar.MenuAction>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content
							class="rounded-lg"
							side={sidebar.isMobile ? 'bottom' : 'right'}
							align={sidebar.isMobile ? 'end' : 'start'}
						>
							<!-- eslint-disable svelte/no-navigation-without-resolve -->
							<DropdownMenu.Item
								onclick={() => {
									goto('#?edit');
								}}
							>
								<PencilIcon class="text-muted-foreground" />
								<span>{m.edit()}</span>
							</DropdownMenu.Item>
							<!-- eslint-enable svelte/no-navigation-without-resolve -->
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{/if}
			</Sidebar.MenuItem>
		{/each}
		{#if hasMoreItems}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton onclick={showMoreItems}>
					<EllipsisIcon class="text-sidebar-foreground/70" />
					{m.more()}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}
	</Sidebar.Menu>
</Sidebar.Group>
