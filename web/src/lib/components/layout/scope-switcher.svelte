<script lang="ts">
	import Icon from '@iconify/svelte';

	import { resolve } from '$app/paths';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { scopeIcon } from '$lib/components/scopes/icon';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';

	import DialogCreateScope from './dialog-create-scope.svelte';

	let {
		active = $bindable<string>(),
		scopes,
		tier,
		onSelect
	}: {
		active: string;
		scopes: Scope[];
		tier: string;
		onSelect: (index: number, home?: boolean) => Promise<void>;
	} = $props();

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

	function toggleDialog() {
		open = !open;
	}
</script>

<svelte:window
	use:shortcut={{
		key: '1',
		ctrl: true,
		callback: async () => await onSelect(0)
	}}
	use:shortcut={{
		key: '2',
		ctrl: true,
		callback: async () => await onSelect(1)
	}}
	use:shortcut={{
		key: '3',
		ctrl: true,
		callback: async () => await onSelect(2)
	}}
	use:shortcut={{
		key: '4',
		ctrl: true,
		callback: async () => await onSelect(3)
	}}
	use:shortcut={{
		key: '5',
		ctrl: true,
		callback: async () => await onSelect(4)
	}}
	use:shortcut={{
		key: '6',
		ctrl: true,
		callback: async () => await onSelect(5)
	}}
	use:shortcut={{
		key: '7',
		ctrl: true,
		callback: async () => await onSelect(6)
	}}
	use:shortcut={{
		key: '8',
		ctrl: true,
		callback: async () => await onSelect(7)
	}}
	use:shortcut={{
		key: '9',
		ctrl: true,
		callback: async () => await onSelect(8)
	}}
/>
<DialogCreateScope bind:open />

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						data-state={props['data-state']}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<Button
							href={resolve('/(auth)/scope/[scope]/setup', { scope: active })}
							class="group/icon flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground transition"
							disabled={active === 'OtterScale'}
						>
							<Icon
								icon="{scopeIcon(scopes.findIndex((s) => s.name === active))}-fill"
								class="size-4.5 group-hover/icon:hidden"
							/>
							<Icon icon="ph:wrench-fill" class="hidden size-4.5 group-hover/icon:block" />
						</Button>
						<div {...props} class="flex h-12 w-full items-center">
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-medium">{active}</span>
								<span class="truncate text-xs">{tier}</span>
							</div>
							<Icon icon="ph:caret-up-down-bold" class="ml-auto size-4 shrink-0" />
						</div>
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">{m.scopes()}</DropdownMenu.Label>
				{#each scopes as scope, index (scope.name)}
					<DropdownMenu.Item onSelect={async () => await onSelect(index, true)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center rounded-md border">
							<Icon icon="{scopeIcon(index)}-bold" class="size-3.5 shrink-0" />
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
					<div class="font-medium text-muted-foreground">{m.add_scope()}</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
