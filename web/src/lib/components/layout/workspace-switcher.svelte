<script lang="ts">
	import ActivityIcon from '@lucide/svelte/icons/activity';
	import ApertureIcon from '@lucide/svelte/icons/aperture';
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import BlocksIcon from '@lucide/svelte/icons/blocks';
	import BoxIcon from '@lucide/svelte/icons/box';
	import BriefcaseIcon from '@lucide/svelte/icons/briefcase';
	import ChartPieIcon from '@lucide/svelte/icons/chart-pie';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import CircleDashedIcon from '@lucide/svelte/icons/circle-dashed';
	import CommandIcon from '@lucide/svelte/icons/command';
	import ComponentIcon from '@lucide/svelte/icons/component';
	import CrownIcon from '@lucide/svelte/icons/crown';
	import CylinderIcon from '@lucide/svelte/icons/cylinder';
	import DiamondIcon from '@lucide/svelte/icons/diamond';
	import DiscIcon from '@lucide/svelte/icons/disc';
	import DnaIcon from '@lucide/svelte/icons/dna';
	import FrameIcon from '@lucide/svelte/icons/frame';
	import GalleryVerticalEndIcon from '@lucide/svelte/icons/gallery-vertical-end';
	import GemIcon from '@lucide/svelte/icons/gem';
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import GridIcon from '@lucide/svelte/icons/grid';
	import HashIcon from '@lucide/svelte/icons/hash';
	import HexagonIcon from '@lucide/svelte/icons/hexagon';
	import LayersIcon from '@lucide/svelte/icons/layers';
	import LayoutGridIcon from '@lucide/svelte/icons/layout-grid';
	import LibraryIcon from '@lucide/svelte/icons/library';
	import MedalIcon from '@lucide/svelte/icons/medal';
	import MountainIcon from '@lucide/svelte/icons/mountain';
	import PackageIcon from '@lucide/svelte/icons/package';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import PyramidIcon from '@lucide/svelte/icons/pyramid';
	import RadarIcon from '@lucide/svelte/icons/radar';
	import ShieldIcon from '@lucide/svelte/icons/shield';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import { type TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import type { Component } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import DialogCreateWorkspace from '$lib/components/layout/dialog-create-workspace.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import { m } from '$lib/paraglide/messages';
	import type { User } from '$lib/server/session';

	let {
		scope,
		workspaces,
		user
	}: { scope: string; workspaces: TenantOtterscaleIoV1Alpha1Workspace[]; user: User } = $props();
	const sidebar = useSidebar();
	let activeWorkspace: TenantOtterscaleIoV1Alpha1Workspace | undefined = $derived(
		workspaces.length > 0 ? workspaces[0] : undefined
	);

	let createWorkspaceOpen = $state(false);

	const workspaceIcons: Component[] = [
		ActivityIcon,
		ApertureIcon,
		AudioWaveformIcon,
		BlocksIcon,
		BoxIcon,
		BriefcaseIcon,
		ChartPieIcon,
		CircleDashedIcon,
		CommandIcon,
		ComponentIcon,
		CrownIcon,
		CylinderIcon,
		DiamondIcon,
		DiscIcon,
		DnaIcon,
		FrameIcon,
		GalleryVerticalEndIcon,
		GemIcon,
		GlobeIcon,
		GridIcon,
		HashIcon,
		HexagonIcon,
		LayersIcon,
		LayoutGridIcon,
		LibraryIcon,
		MedalIcon,
		MountainIcon,
		PackageIcon,
		PyramidIcon,
		RadarIcon,
		ShieldIcon,
		ZapIcon
	];

	function getWorkspaceIcon(name: string | undefined): Component {
		if (!name) return PlusIcon;

		let hash = 0;
		for (let i = 0; i < name.length; i++) {
			hash = name.charCodeAt(i) + ((hash << 5) - hash);
		}

		const index = Math.abs(hash % workspaceIcons.length);
		return workspaceIcons[index];
	}

	let ActiveIcon = $derived(getWorkspaceIcon(activeWorkspace?.metadata?.name));

	function onSelect(index: number) {
		if (index >= workspaces.length) {
			return;
		}
		activeWorkspace = workspaces[index];

		goto(
			resolve(
				`/(auth)/${scope}/Workspace/workspaces?group=tenant.otterscale.io&version=v1alpha1&name=${activeWorkspace?.metadata?.name ?? ''}`
			)
		);

		// TODO: Add logic to update global workspace state or navigate to the selected workspace's route.
		// await goto(resolve('/(auth)/scope/[scope]', { scope: scope.name }));
		if (activeWorkspace.metadata?.name) {
			toast.success(m.switch_workspace({ name: activeWorkspace.metadata.name }));
		}
	}
</script>

<svelte:window
	use:shortcut={{
		key: '1',
		ctrl: true,
		callback: () => onSelect(0)
	}}
	use:shortcut={{
		key: '2',
		ctrl: true,
		callback: () => onSelect(1)
	}}
	use:shortcut={{
		key: '3',
		ctrl: true,
		callback: () => onSelect(2)
	}}
	use:shortcut={{
		key: '4',
		ctrl: true,
		callback: () => onSelect(3)
	}}
	use:shortcut={{
		key: '5',
		ctrl: true,
		callback: () => onSelect(4)
	}}
	use:shortcut={{
		key: '6',
		ctrl: true,
		callback: () => onSelect(5)
	}}
	use:shortcut={{
		key: '7',
		ctrl: true,
		callback: () => onSelect(6)
	}}
	use:shortcut={{
		key: '8',
		ctrl: true,
		callback: () => onSelect(7)
	}}
	use:shortcut={{
		key: '9',
		ctrl: true,
		callback: () => onSelect(8)
	}}
/>
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
							class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
						>
							<ActiveIcon class="size-4" />
							{#if !activeWorkspace}
								<span class="absolute top-1 left-8 flex size-3">
									<span
										class="absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75"
									></span>
									<span class="relative inline-flex size-3 rounded-full bg-blue-500"></span>
								</span>
							{/if}
						</div>
						<div class="grid flex-1 text-start text-sm leading-tight">
							{#if activeWorkspace}
								<span class="truncate font-medium"> {activeWorkspace.metadata?.name} </span>
								<span class="flex items-center gap-2 truncate text-xs text-muted-foreground">
									{activeWorkspace.spec.users.find((u) => u.subject === user.sub)?.role}
								</span>
							{:else}
								<span class="truncate font-medium"> OtterScale </span>
								<span class="flex items-center gap-2 truncate text-xs text-muted-foreground">
									{m.no_workspace_selected()}
								</span>
							{/if}
						</div>
						<ChevronsUpDownIcon class="ms-auto" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">
					{m.workspace()}
				</DropdownMenu.Label>
				{#each workspaces as workspace, index (workspace.metadata?.name)}
					{@const WorkspaceIcon = getWorkspaceIcon(workspace.metadata?.name)}
					<DropdownMenu.Item onSelect={() => onSelect(index)} class="gap-2 p-2">
						<div class="flex size-6 items-center justify-center rounded-md border">
							<WorkspaceIcon class="size-3 shrink-0" />
						</div>
						<div class="grid flex-1 text-start text-xs leading-tight">
							<span class="truncate font-medium">{workspace.metadata?.name}</span>
							<span class="truncate text-[10px] text-muted-foreground"
								>{workspace.spec.users.find((u) => u.subject === user.sub)?.role}</span
							>
						</div>
						<DropdownMenu.Shortcut class="flex items-center gap-0.5 text-sm">
							<CommandIcon class="size-3" />
							{#if index < 9}
								<span class="font-mono">{index + 1}</span>
							{/if}
						</DropdownMenu.Shortcut>
					</DropdownMenu.Item>
				{/each}
				<DropdownMenu.Separator />
				<DropdownMenu.Item class="gap-2 p-2" onSelect={() => (createWorkspaceOpen = true)}>
					<div class="flex size-6 items-center justify-center rounded-md border bg-transparent">
						<PlusIcon class="size-3.5" />
					</div>
					<div class="text-xs font-medium text-muted-foreground">{m.add_workspace()}</div>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>

<DialogCreateWorkspace bind:open={createWorkspaceOpen} />
