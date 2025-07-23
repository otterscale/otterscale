<script lang="ts" module>
	import {
		applicationsPath,
		applicationsServicePath,
		applicationsStorePath,
		applicationsWorkloadPath,
		documentationPath,
		feedbackPath,
		machinesMetalPath,
		machinesPath,
		machinesVirtualMachinePath,
		modelLLMPath,
		modelPath,
		settingsPath,
		settingsNetworkPath,
		storageBlockDevicePath,
		storageClusterPath,
		storageFileSystemPath,
		storageObjectGatewayPath,
		storagePath
	} from '$lib/path';

	const NAVIGATION_DATA = {
		main: [
			{
				title: 'Model',
				url: modelPath,
				items: [{ title: 'LLM', url: modelLLMPath }]
			},
			{
				title: 'Applications',
				url: applicationsPath,
				items: [
					{ title: 'Workload', url: applicationsWorkloadPath },
					{ title: 'Service', url: applicationsServicePath },
					{ title: 'Store', url: applicationsStorePath }
				]
			},
			{
				title: 'Storage',
				url: storagePath,
				items: [
					{ title: 'Cluster', url: storageClusterPath },
					{ title: 'Block Device', url: storageBlockDevicePath },
					{ title: 'File System - NFS', url: storageFileSystemPath },
					{ title: 'Object Gateway - S3', url: storageObjectGatewayPath }
				]
			},
			{
				title: 'Machines',
				url: machinesPath,
				items: [
					{ title: 'Metal', url: machinesMetalPath },
					{ title: 'Virtual Machine', url: machinesVirtualMachinePath }
				]
			},
			{
				title: 'Settings',
				url: settingsPath,
				items: [
					{ title: 'General', url: settingsPath },
					{ title: 'Network', url: settingsNetworkPath }
				]
			}
		],
		secondary: [
			{ title: 'Feedback', url: feedbackPath },
			{ title: 'Documentation', url: documentationPath }
		],
		bookmarks: [
			{ name: 'FOO 1', url: '#' },
			{ name: 'BAR 1', url: '#' },
			{ name: 'FOO 2', url: '#' },
			{ name: 'BAR 2', url: '#' },
			{ name: 'FOO 3', url: '#' },
			{ name: 'BAR 3', url: '#' }
		]
	};
</script>

<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import NavMain from './nav-main.svelte';
	import NavPrimary from './nav-primary.svelte';
	import NavSecondary from './nav-secondary.svelte';
	import NavUser from './nav-user.svelte';
	import ScopeSwitcher from './scope-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import type { ComponentProps } from 'svelte';
	import type { User } from 'better-auth';

	type Props = { user: User } & ComponentProps<typeof Sidebar.Root>;

	let { user, ref = $bindable(null), ...restProps }: Props = $props();

	let mounted = $state(false);
	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const scopes = writable<Scope[]>([]);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	onMount(async () => {
		await fetchScopes();
		mounted = true;
	});
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		{#if mounted}
			<ScopeSwitcher scopes={$scopes} />
		{:else}
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg">
						{#snippet child({ props })}
							<div {...props}>
								<Skeleton class="bg-sidebar-primary/50 size-8 rounded-lg" />
								<div class="grid flex-1 space-y-1 text-left text-sm leading-tight">
									<Skeleton class="bg-sidebar-primary/50 h-3 w-[150px]" />
									<Skeleton class="bg-sidebar-primary/50 h-3 w-[50px]" />
								</div>
							</div>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{/if}
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={NAVIGATION_DATA.main} />
		<NavPrimary bookmarks={NAVIGATION_DATA.bookmarks} />
		<NavSecondary items={NAVIGATION_DATA.secondary} class="mt-auto" />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
