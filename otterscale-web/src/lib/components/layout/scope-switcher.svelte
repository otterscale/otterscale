<script lang="ts">
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { activeScope } from '$lib/stores';
	import DialogCreateScope from './dialog-create-scope.svelte';
	import { m } from '$lib/paraglide/messages';

	let { scopes, edition }: { scopes: Scope[]; edition: string } = $props();
	let open = $state(false);

	const sidebar = useSidebar();

	const SHORTCUT_ICONS = [
		'ph:number-one',
		'ph:number-two',
		'ph:number-three',
		'ph:number-four',
		'ph:number-five',
		'ph:number-six',
		'ph:number-seven',
		'ph:number-eight',
		'ph:number-nine'
	];

	const SCOPE_ICONS = [
		'ph:airplane-tilt',
		'ph:cactus',
		'ph:cherries',
		'ph:piggy-bank',
		'ph:flower',
		'ph:joystick',
		'ph:clover',
		'ph:cube',
		'ph:gavel'
	];

	function getIcon(name: string): string {
		const index = scopes.findIndex((scope) => scope.name === name);
		return index !== -1 ? SCOPE_ICONS[index % SCOPE_ICONS.length] : SCOPE_ICONS[0];
	}

	function handleScopeShortcut(index: number): void {
		if (scopes.length > index) {
			activeScope.set(scopes[index]);
			toast.info(`Switched to '${$activeScope.name}' scope`);
		}
	}

	function toggleDialog(): void {
		open = !open;
	}
</script>

<svelte:window
	use:shortcut={{
		key: '1',
		ctrl: true,
		callback: () => handleScopeShortcut(0)
	}}
	use:shortcut={{
		key: '2',
		ctrl: true,
		callback: () => handleScopeShortcut(1)
	}}
	use:shortcut={{
		key: '3',
		ctrl: true,
		callback: () => handleScopeShortcut(2)
	}}
	use:shortcut={{
		key: '4',
		ctrl: true,
		callback: () => handleScopeShortcut(3)
	}}
	use:shortcut={{
		key: '5',
		ctrl: true,
		callback: () => handleScopeShortcut(4)
	}}
	use:shortcut={{
		key: '6',
		ctrl: true,
		callback: () => handleScopeShortcut(5)
	}}
	use:shortcut={{
		key: '7',
		ctrl: true,
		callback: () => handleScopeShortcut(6)
	}}
	use:shortcut={{
		key: '8',
		ctrl: true,
		callback: () => handleScopeShortcut(7)
	}}
	use:shortcut={{
		key: '9',
		ctrl: true,
		callback: () => handleScopeShortcut(8)
	}}
/>

<DialogCreateScope bind:open />

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
							class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
						>
							<Icon icon="{getIcon($activeScope.name)}-fill" class="size-4.5" />
						</div>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{$activeScope.name}</span>
							<span class="truncate text-xs">{edition}</span>
						</div>
						<Icon icon="ph:caret-up-down-bold" class="ml-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-muted-foreground text-xs">{m.scopes()}</DropdownMenu.Label>
				{#each scopes as scope, index (scope.name)}
					<DropdownMenu.Item onSelect={() => handleScopeShortcut(index)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center rounded-md border">
							<Icon icon="{getIcon(scope.name)}-bold" class="size-3.5 shrink-0" />
						</div>
						{scope.name}
						{#if index < 9}
							<DropdownMenu.Shortcut>
								<div class="flex items-center justify-center">
									<Icon icon="ph:control-bold" />
									<Icon icon={SHORTCUT_ICONS[index]} class="size-3.5" />
								</div>
							</DropdownMenu.Shortcut>
						{/if}
					</DropdownMenu.Item>
				{/each}
				<DropdownMenu.Separator />
				<DropdownMenu.Item class="gap-2 p-2" onclick={toggleDialog}>
					<div class="flex size-6 items-center justify-center rounded-md border bg-transparent">
						<Icon icon="ph:plus" />
					</div>
					<div class="text-muted-foreground font-medium">{m.add_scope()}</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
