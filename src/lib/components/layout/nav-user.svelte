<script lang="ts">
	import BellRingIcon from '@lucide/svelte/icons/bell-ring';
	import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import MoonIcon from '@lucide/svelte/icons/moon';
	import SunIcon from '@lucide/svelte/icons/sun';
	import UserRoundIcon from '@lucide/svelte/icons/user-round';
	import { mode, toggleMode } from 'mode-watcher';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import { m } from '$lib/paraglide/messages';
	import { getLocale, type Locale, setLocale } from '$lib/paraglide/runtime';
	import type { User } from '$lib/server/session';

	import SheetNotification from './sheet-notification.svelte';

	let { user }: { user: User } = $props();

	let locale = $state(getLocale());
	let open = $state(false);

	const sidebar = useSidebar();

	const getUserInitials = (name: string) => {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase();
	};

	const handleSignOut = () => {
		goto(resolve('/logout'));
	};

	const handleLanguageChange = (newLocale: Locale) => {
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
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						{...props}
					>
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user.picture} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{getUserInitials(user.name)}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{user.name}</span>
							<span class="truncate text-xs text-muted-foreground">{user.email}</span>
						</div>
						<ChevronsUpDownIcon class="ms-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user.picture} alt={user.name} />
							<Avatar.Fallback class="rounded-lg">{getUserInitials(user.name)}</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-start text-sm leading-tight">
							<span class="truncate font-medium">{user.name}</span>
							<span class="truncate text-xs text-muted-foreground">{user.email}</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={() => goto(resolve('/(auth)/account'))}>
						<UserRoundIcon />
						{m.account()}
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={toggleNotification}>
						<BellRingIcon />
						{m.notifications()}
						<span class="relative flex size-2.5">
							<span
								class="absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75"
							></span>
							<span class="relative inline-flex size-2.5 rounded-full bg-blue-500"></span>
						</span>
						<DropdownMenu.Shortcut>
							<div class="flex items-center justify-center space-x-1">
								<ChevronUpIcon />
								<span>/</span>
							</div>
						</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item onclick={toggleMode}>
						{#if mode.current === 'light'}
							<MoonIcon />
							{m.dark_mode()}
						{:else}
							<SunIcon />
							{m.light_mode()}
						{/if}
					</DropdownMenu.Item>
					<DropdownMenu.Sub>
						<DropdownMenu.SubTrigger>
							<GlobeIcon />
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
				<DropdownMenu.Item variant="destructive" onclick={handleSignOut}>
					<LogOutIcon />
					{m.log_out()}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
