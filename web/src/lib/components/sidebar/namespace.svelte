<script lang="ts">
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';

	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';

	let {
		namespaces
	}: { namespaces: { name: string; plan: string; icon: string; color: string }[] } = $props();
	const sidebar = useSidebar();

	let active = $state(namespaces[0]);
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<div
							class="relative flex aspect-square size-8 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
						>
							<Icon icon={active.icon} class="size-6" />
						</div>
						<div class="grid flex-1 pl-1 text-left text-sm leading-tight">
							<span class="truncate font-semibold">
								{active.name}
							</span>
							<span class="truncate text-xs">{active.plan}</span>
						</div>
						<Icon icon="ph:caret-up-down" class="ml-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-[--bits-dropdown-menu-anchor-width] min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">Namespaces</DropdownMenu.Label>
				{#each namespaces as namespace, index (namespace.name)}
					<DropdownMenu.Item onSelect={() => (active = namespace)} class="gap-2 p-2 text-sm">
						<div class="flex size-7 items-center justify-center rounded-md border">
							<Icon icon={namespace.icon} class="size-5 shrink-0" />
						</div>
						{namespace.name}
						<DropdownMenu.Shortcut>⌘{index + 1}</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				{/each}
				<DropdownMenu.Separator />
				<DropdownMenu.Item class="gap-2 p-2 text-sm">
					<div class="flex size-7 items-center justify-center rounded-md border">
						<Icon icon="ph-plus" class="size-5 shrink-0" />
					</div>
					<Button
						variant="ghost"
						class="h-auto p-0 font-medium text-muted-foreground"
						onclick={() => toast.info('Coming in next version')}
					>
						Add namespace
					</Button>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
