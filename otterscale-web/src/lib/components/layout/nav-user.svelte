<script lang="ts">
	import type { User } from 'better-auth';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { authClient } from '$lib/auth-client';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { accountPath, loginPath, settingsPath } from '$lib/path';
	import SheetBookmark from './sheet-bookmark.svelte';
	import SheetNotification from './sheet-notification.svelte';

	let { user }: { user: User } = $props();
	const sidebar = useSidebar();

	let openBookmark = $state(false);
	let openNotification = $state(false);

	const menuItems = [
		{ icon: 'ph:user-bold', label: 'Account', action: () => goto(accountPath) },
		{ icon: 'ph:push-pin-bold', label: 'Bookmarks', action: () => (openBookmark = true) },
		{
			icon: 'ph:bell-ringing-bold',
			label: 'Notifications',
			action: () => (openNotification = true),
			hasNotification: true
		}
	];

	const handleSignOut = () => {
		authClient.signOut({
			fetchOptions: {
				onSuccess: () => {
					toast.success('Signed out successfully!');
					goto(loginPath);
				}
			}
		});
	};

	const getUserInitials = (name: string) => {
		return (
			name
				?.split(' ')
				.map((n) => n[0])
				.join('')
				.toUpperCase() || ''
		);
	};
</script>

<SheetBookmark bind:open={openBookmark} />
<SheetNotification bind:open={openNotification} />

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
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user.image} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{getUserInitials(user.name)}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{user.name}</span>
							<span class="truncate text-xs">{user.email}</span>
						</div>
						<Icon icon="ph:caret-up-down-bold" class="ml-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>

			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<!-- User Info Header -->
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user.image} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{getUserInitials(user.name)}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{user.name}</span>
							<span class="truncate text-xs">{user.email}</span>
						</div>
					</div>
				</DropdownMenu.Label>

				<DropdownMenu.Separator />

				<!-- Main Menu Items -->
				<DropdownMenu.Group>
					{#each menuItems as item}
						<DropdownMenu.Item onclick={item.action}>
							<Icon icon={item.icon} />
							{item.label}
							{#if item.hasNotification}
								<span class="absolute right-2 flex size-2.5">
									<span
										class="absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75"
									></span>
									<span class="relative inline-flex size-2.5 rounded-full bg-blue-500"></span>
								</span>
							{/if}
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Group>

				<DropdownMenu.Separator />

				<!-- Settings -->
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={() => goto(settingsPath)}>
						<Icon icon="ph:gear-bold" />
						Settings
					</DropdownMenu.Item>
				</DropdownMenu.Group>

				<DropdownMenu.Separator />

				<!-- Sign Out -->
				<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
					<Icon icon="ph:sign-out-bold" />
					Log out
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
