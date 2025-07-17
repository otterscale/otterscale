<script lang="ts">
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar';
	import DialogCreateScope from './dialog-create-scope.svelte';

	let { teams }: { teams: { name: string; icon: any; enterprise: boolean }[] } = $props();
	let activeTeam = $state(teams[0]);
	let open = $state(false);

	const sidebar = useSidebar();

	const indexIcons = [
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

	function handleTeamShortcut(index: number) {
		if (teams.length > index) {
			activeTeam = teams[index];
			toast.info(`Toggle to ${activeTeam.name}`);
		}
	}

	const toggleDialog = () => {
		open = !open;
	};
</script>

<svelte:window
	use:shortcut={{
		key: '1',
		ctrl: true,
		callback: () => handleTeamShortcut(0)
	}}
	use:shortcut={{
		key: '2',
		ctrl: true,
		callback: () => handleTeamShortcut(1)
	}}
	use:shortcut={{
		key: '3',
		ctrl: true,
		callback: () => handleTeamShortcut(2)
	}}
	use:shortcut={{
		key: '4',
		ctrl: true,
		callback: () => handleTeamShortcut(3)
	}}
	use:shortcut={{
		key: '5',
		ctrl: true,
		callback: () => handleTeamShortcut(4)
	}}
	use:shortcut={{
		key: '6',
		ctrl: true,
		callback: () => handleTeamShortcut(5)
	}}
	use:shortcut={{
		key: '7',
		ctrl: true,
		callback: () => handleTeamShortcut(6)
	}}
	use:shortcut={{
		key: '8',
		ctrl: true,
		callback: () => handleTeamShortcut(7)
	}}
	use:shortcut={{
		key: '9',
		ctrl: true,
		callback: () => handleTeamShortcut(8)
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
							<Icon icon={activeTeam.icon + '-fill'} class="size-4.5" />
						</div>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{activeTeam.name}</span>
							<span class="truncate text-xs">{activeTeam.enterprise ? 'Enterprise' : 'Free'}</span>
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
				<DropdownMenu.Label class="text-muted-foreground text-xs">Scopes</DropdownMenu.Label>
				{#each teams as team, index (team.name)}
					<DropdownMenu.Item onSelect={() => handleTeamShortcut(index)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center rounded-md border">
							<Icon icon={team.icon + '-bold'} class="size-3.5 shrink-0" />
						</div>
						{team.name}
						{#if index < 9}
							<DropdownMenu.Shortcut>
								<div class="flex items-center justify-center">
									<Icon icon="ph:control-bold" />
									<Icon icon={indexIcons[index]} />
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
					<div class="text-muted-foreground font-medium">Add scope</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
