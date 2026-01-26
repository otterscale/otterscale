<script lang="ts">
	import 'driver.js/dist/driver.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import CombineIcon from '@lucide/svelte/icons/combine';
	import HelpCircleIcon from '@lucide/svelte/icons/help-circle';
	import HouseIcon from '@lucide/svelte/icons/house';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import { getContext, onMount, type Snippet } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { type Scope, ScopeService } from '$lib/api/scope/v1/scope_pb';
	import {
		navData,
		NavGeneral,
		NavMain,
		NavOverview,
		NavSecondary,
		NavUser,
		startTour,
		WorkspaceSwitcher
	} from '$lib/components/layout';
	import { globalRoutes, platformRoutes } from '$lib/components/layout/routes';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	import type { LayoutData } from './$types';

	let {
		data,
		children
	}: {
		data: LayoutData;
		children: Snippet;
	} = $props();

	const current = $derived($breadcrumbs.at(-1));

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const resourceClient = createClient(ResourceService, transport);

	let activeScope = $state(page.params.scope ?? '');
	let scopes = $state<Scope[]>([]);
	let workspaces = $state<TenantOtterscaleIoV1Alpha1Workspace[]>([]);
	let next = $state(false);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes = response.scopes.filter((scope) => scope.name !== 'cos');
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function fetchWorkspaces(cluster: string) {
		try {
			const response = await resourceClient.list({
				cluster: cluster,
				group: 'tenant.otterscale.io',
				version: 'v1alpha1',
				resource: 'workspaces',
				labelSelector: 'user.otterscale.io/' + data.user.sub
			});
			workspaces = response.items.map((item) => item.object as TenantOtterscaleIoV1Alpha1Workspace);
		} catch (error) {
			console.error('Failed to fetch workspaces:', error);
		}
	}

	async function onValueChange(cluster: string) {
		await fetchWorkspaces(cluster);
		await goto(resolve('/(auth)/scope/[scope]', { scope: cluster }));
		toast.success(m.switch_scope({ name: cluster }));
	}

	async function onHomeClick() {
		activeScope = '';
		await goto(resolve('/(auth)/console'));
	}

	let isMounted = $state(false);
	onMount(async () => {
		await fetchScopes();

		if (activeScope) {
			await fetchWorkspaces(activeScope);
		}

		isMounted = true;
	});
</script>

<svelte:head>
	<title>{current ? `${current.title} - OtterScale` : 'OtterScale'}</title>
</svelte:head>

<Sidebar.Provider>
	<Sidebar.Root id="sidebar-guide-step" collapsible="icon" variant="inset" class="p-3">
		{#if activeScope && isMounted}
			<Sidebar.Header id="workspace-guide-step">
				<WorkspaceSwitcher {workspaces} user={data.user} />
			</Sidebar.Header>
			<Sidebar.Content class="gap-2">
				<NavOverview items={navData.overview} />
				{#if next}
					<NavMain label="AI Studio" items={navData.aiStudio} />
					<NavMain label="Applications" items={navData.applications} />
					<NavMain label="Resources" items={navData.resources} />
					<NavMain label="Governance" items={navData.governance} />
					<NavMain label="Reliability" items={navData.reliability} />
					<NavMain label="System" items={navData.system} />
				{:else}
					<NavGeneral title={m.platform()} routes={platformRoutes(activeScope)} />
					<NavGeneral title={m.global()} routes={globalRoutes()} />
				{/if}
			</Sidebar.Content>
			<Button
				class="mx-auto w-full text-xs text-muted-foreground"
				variant="link"
				onclick={() => (next = !next)}
			>
				{#if next}
					<ChevronLeftIcon class="size-3.5" />
					{m.switch_to_classic()}
				{:else}
					<ZapIcon class="size-3.5" />
					{m.try_next_version()}
				{/if}
			</Button>
		{:else}
			<Sidebar.Header id="workspace-guide-step">
				<div class="flex h-12 w-full items-center gap-2 overflow-hidden rounded-md p-2">
					<Skeleton class="size-8 bg-foreground/10" />
					<div class="space-y-2">
						<Skeleton class="h-3 w-36 bg-foreground/10" />
						<Skeleton class="h-2 w-12 bg-foreground/10" />
					</div>
				</div>
			</Sidebar.Header>
			<Sidebar.Content class="gap-2">
				<div class="relative flex w-full min-w-0 flex-col space-y-4 px-4 py-2">
					<Skeleton class="h-3 w-8 bg-foreground/10" />
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
				</div>
				<div class="relative flex w-full min-w-0 flex-col space-y-4 px-4 py-2">
					<Skeleton class="h-3 w-8 bg-foreground/10" />
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
					<div class="flex items-center space-x-2">
						<Skeleton class="h-4 w-4 bg-foreground/10" />
						<Skeleton class="h-4 w-32 bg-foreground/10" />
					</div>
				</div>
			</Sidebar.Content>
		{/if}
		<NavSecondary />
		<Sidebar.Footer>
			<NavUser user={data.user} />
		</Sidebar.Footer>
		<Sidebar.Rail />
	</Sidebar.Root>
	<Sidebar.Inset>
		<header
			class="flex h-16 shrink-0 items-center justify-between gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
		>
			<div class="flex items-center gap-2 px-4">
				<Sidebar.Trigger class="-ms-1" />
				<Separator orientation="vertical" class="me-2 data-[orientation=vertical]:h-4" />
				<Breadcrumb.Root>
					<Breadcrumb.List>
						{#each $breadcrumbs as item (item.url)}
							{#if item.url === current?.url}
								<Breadcrumb.Item>
									<Breadcrumb.Page>{current.title}</Breadcrumb.Page>
								</Breadcrumb.Item>
							{:else}
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href={item.url}>{item.title}</Breadcrumb.Link>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
							{/if}
						{/each}
					</Breadcrumb.List>
				</Breadcrumb.Root>
			</div>
			<div class="flex items-center gap-2 px-4">
				<Button variant="ghost" size="icon" class="size-7" onclick={startTour}>
					<HelpCircleIcon />
					<span class="sr-only">Help</span>
				</Button>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button {...props} id="cluster-guide-step" variant="ghost" size="icon" class="size-7">
								<CombineIcon />
								<span class="sr-only">Toggle Clusters</span>
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-40" align="end">
						<DropdownMenu.Group>
							<DropdownMenu.Label>{m.cluster()}</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.RadioGroup bind:value={activeScope} {onValueChange}>
								{#each scopes as scope, index (index)}
									<DropdownMenu.RadioItem value={scope.name}>{scope.name}</DropdownMenu.RadioItem>
								{/each}
							</DropdownMenu.RadioGroup>
						</DropdownMenu.Group>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
				<Button variant="ghost" size="icon" class="size-7" onclick={onHomeClick}>
					<HouseIcon />
					<span class="sr-only">Back to Home</span>
				</Button>
			</div>
		</header>
		<main class="flex flex-1 flex-col px-2 md:px-4 lg:px-8">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
