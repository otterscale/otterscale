<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import * as m from '$lib/paraglide/messages.js';
	import { i18n } from '$lib/i18n';
	import { goto } from '$app/navigation';

	let { user }: { user: { name: string; email: string; avatar: string; fallback: string } } =
		$props();
	const sidebar = useSidebar();
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						{...props}
					>
						<Avatar.Root class="h-8 w-8 rounded-lg">
							<Avatar.Image src={user.avatar} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{user.fallback}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-semibold">{user.name}</span>
							<span class="truncate text-xs">{user.email}</span>
						</div>
						<Icon icon="ph:caret-up-down" class="ml-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-[--bits-dropdown-menu-anchor-width] min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
						<Avatar.Root class="h-8 w-8 rounded-lg">
							<Avatar.Image src={user.avatar} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{user.fallback}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-semibold">{user.name}</span>
							<span class="truncate text-xs">{user.email}</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings#profile'))}>
						<Icon icon="ph:user" />
						<span class="pl-2">{m.avatar_profile()}</span>
						<DropdownMenu.Shortcut>⇧⌘P</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings'))}>
						<Icon icon="ph:gear" />
						<span class="pl-2">{m.avatar_settings()}</span>
						<DropdownMenu.Shortcut>⇧⌘S</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/settings#superuser'))}>
						<Icon icon="ph:key" />
						<span class="pl-2">{m.avatar_superuser()}</span>
						<DropdownMenu.Shortcut>⇧⌘U</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/favorites'))}>
						<Icon icon="ph:clover" />
						<span class="pl-2">{m.avatar_favorites()}</span>
						<DropdownMenu.Shortcut>⇧⌘F</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/recents'))}>
						<Icon icon="ph:clock" />
						<span class="pl-2">{m.avatar_recents()}</span>
						<DropdownMenu.Shortcut>⇧⌘R</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item onclick={() => goto(i18n.resolveRoute('/logout'))}>
					<Icon icon="ph:sign-in" />
					<span class="pl-2">{m.avatar_sign_out()}</span>
					<DropdownMenu.Shortcut>⇧⌘Q</DropdownMenu.Shortcut>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
