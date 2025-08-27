<script lang="ts">
	import type { User } from 'better-auth';
	import { mode, toggleMode } from 'mode-watcher';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import { authClient } from '$lib/auth-client';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { getLocale, setLocale } from '$lib/paraglide/runtime';
	import { m } from '$lib/paraglide/messages.js';
	import { dynamicPaths, staticPaths } from '$lib/path';
	import SheetNotification from './sheet-notification.svelte';

	let { user }: { user: User } = $props();
	let locale = $state(getLocale());
	let open = $state(false);

	const sidebar = useSidebar();

	const getUserInitials = (name: string) => {
		return (
			name
				?.split(' ')
				.map((n) => n[0])
				.join('')
				.toUpperCase() || ''
		);
	};

	const handleSignOut = () => {
		authClient.signOut({
			fetchOptions: {
				onSuccess: () => {
					toast.success(m.sign_out_success());
					goto(staticPaths.login.url);
				}
			}
		});
	};

	const handleLanguageChange = (newLocale: any) => {
		setLocale(newLocale);
		locale = newLocale;
	};

	const toggleNotification = () => {
		open = !open;
	};
</script>

<svelte:window
	use:shortcut={{
		key: '/',
		ctrl: true,
		callback: toggleNotification
	}}
/>

<SheetNotification bind:open />

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

				<!-- User Actions -->
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={() => goto(dynamicPaths.account(page.params.scope).url)}>
						<Icon icon="ph:user-bold" />
						{m.account()}
					</DropdownMenu.Item>

					<DropdownMenu.Item onclick={toggleNotification}>
						<Icon icon="ph:bell-ringing-bold" />
						{m.notifications()}
						<span class="relative flex size-2.5">
							<span
								class="absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75"
							></span>
							<span class="relative inline-flex size-2.5 rounded-full bg-blue-500"></span>
						</span>
						<DropdownMenu.Shortcut>
							<div class="flex items-center justify-center space-x-1">
								<Icon icon="ph:control-bold" />
								<span>/</span>
							</div>
						</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				</DropdownMenu.Group>

				<DropdownMenu.Separator />

				<!-- Preferences -->
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={toggleMode}>
						<Icon icon={mode.current === 'light' ? 'ph:sun' : 'ph:moon'} />
						{mode.current === 'light' ? m.light_mode() : m.dark_mode()}
					</DropdownMenu.Item>

					<DropdownMenu.Sub>
						<DropdownMenu.SubTrigger>
							<Icon icon="ph:globe" />
							{m.locale()}
						</DropdownMenu.SubTrigger>
						<DropdownMenu.SubContent>
							<DropdownMenu.RadioGroup bind:value={locale}>
								<DropdownMenu.RadioItem value="en" onclick={() => handleLanguageChange('en')}>
									English
								</DropdownMenu.RadioItem>
								<DropdownMenu.RadioItem
									value="zh-hant"
									onclick={() => handleLanguageChange('zh-hant')}
								>
									繁體中文
								</DropdownMenu.RadioItem>
							</DropdownMenu.RadioGroup>
						</DropdownMenu.SubContent>
					</DropdownMenu.Sub>
				</DropdownMenu.Group>

				<DropdownMenu.Separator />

				<!-- Sign Out -->
				<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
					<Icon icon="ph:sign-out-bold" />
					{m.log_out()}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
