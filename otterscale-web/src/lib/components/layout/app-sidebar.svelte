<script lang="ts" module>
	import { m } from '$lib/paraglide/messages.js';
	import {
		applicationsPath,
		applicationsServicePath,
		applicationsStorePath,
		applicationsWorkloadPath,
		machinesMetalPath,
		machinesPath,
		machinesVirtualMachinePath,
		modelsLLMPath,
		modelsPath,
		settingsPath,
		settingsNetworkPath,
		storageBlockDevicePath,
		storageClusterPath,
		storageFileSystemPath,
		storageObjectGatewayPath,
		storagePath,
		settingsBISTPath,
		databasesPath,
		databasesRelationalPath,
		databasesNoSQLPath
	} from '$lib/path';

	const NAVIGATION_DATA = {
		main: [
			{
				title: m.models(),
				url: modelsPath,
				items: [{ title: m.llm(), url: modelsLLMPath }]
			},
			{
				title: m.databases(),
				url: databasesPath,
				items: [
					{ title: m.relational(), url: databasesRelationalPath },
					{ title: m.no_sql(), url: databasesNoSQLPath }
				]
			},
			{
				title: m.applications(),
				url: applicationsPath,
				items: [
					{ title: m.workload(), url: applicationsWorkloadPath },
					{ title: m.service(), url: applicationsServicePath },
					{ title: m.store(), url: applicationsStorePath }
				]
			},
			{
				title: m.storage(),
				url: storagePath,
				items: [
					{ title: m.cluster(), url: storageClusterPath },
					{ title: m.block_device(), url: storageBlockDevicePath },
					{ title: m.file_system_nfs(), url: storageFileSystemPath },
					{ title: m.object_gateway_s3(), url: storageObjectGatewayPath }
				]
			},
			{
				title: m.machines(),
				url: machinesPath,
				items: [
					{ title: m.metal(), url: machinesMetalPath },
					{ title: m.virtual_machine(), url: machinesVirtualMachinePath }
				]
			},
			{
				title: m.settings(),
				url: settingsPath,
				items: [
					{ title: m.general(), url: settingsPath },
					{ title: m.network(), url: settingsNetworkPath },
					{ title: m.built_in_test(), url: settingsBISTPath }
				]
			}
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
	import type { ComponentProps } from 'svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import type { User } from 'better-auth';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { activeScope, scopeLoading } from '$lib/stores';
	import NavMain from './nav-main.svelte';
	import NavPrimary from './nav-primary.svelte';
	import NavSecondary from './nav-secondary.svelte';
	import NavUser from './nav-user.svelte';
	import ScopeSwitcher from './scope-switcher.svelte';

	type Props = { user: User } & ComponentProps<typeof Sidebar.Root>;

	let { user, ref = $bindable(null), ...restProps }: Props = $props();

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const scopes = writable<Scope[]>([]);

	async function initializeScopes() {
		scopeLoading.set(true);

		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);

			if (response.scopes.length > 0) {
				activeScope.set(response.scopes[0]);
			}
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		} finally {
			scopeLoading.set(false);
		}
	}

	onMount(initializeScopes);

	function renderLoadingSkeleton() {
		return {
			avatar: 'bg-sidebar-primary/50 size-8 rounded-lg',
			title: 'bg-sidebar-primary/50 h-3 w-[150px]',
			subtitle: 'bg-sidebar-primary/50 h-3 w-[50px]'
		};
	}

	const skeletonClasses = renderLoadingSkeleton();
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		{#if $scopeLoading}
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg">
						{#snippet child({ props })}
							<div {...props}>
								<Skeleton class={skeletonClasses.avatar} />
								<div class="grid flex-1 space-y-1 text-left text-sm leading-tight">
									<Skeleton class={skeletonClasses.title} />
									<Skeleton class={skeletonClasses.subtitle} />
								</div>
							</div>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{:else}
			<ScopeSwitcher scopes={$scopes} />
		{/if}
	</Sidebar.Header>

	<Sidebar.Content>
		<NavMain items={NAVIGATION_DATA.main} />
		<NavPrimary bookmarks={NAVIGATION_DATA.bookmarks} />
		<NavSecondary class="mt-auto" />
	</Sidebar.Content>

	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
