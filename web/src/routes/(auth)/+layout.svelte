<script lang="ts">
	import { toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import HouseIcon from '@lucide/svelte/icons/house';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import type { Snippet } from 'svelte';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import {
		type APIResource,
		type DiscoveryRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';
	import { type Scope, ScopeService } from '$lib/api/scope/v1/scope_pb';
	import NavBreadcrumb from '$lib/components/layout/nav-breadcrumb.svelte';
	import { navData } from '$lib/components/layout/nav-data';
	import NavGeneral from '$lib/components/layout/nav-general.svelte';
	import NavMain from '$lib/components/layout/nav-main.svelte';
	import NavDashboard from '$lib/components/layout/nav-overview.svelte';
	import NavSecondary from '$lib/components/layout/nav-secondary.svelte';
	import NavUser from '$lib/components/layout/nav-user.svelte';
	import { globalRoutes, platformRoutes } from '$lib/components/layout/routes';
	import ScopeSwitcher from '$lib/components/layout/scope-switcher.svelte';
	import WorkspaceSwitcher from '$lib/components/layout/workspace-switcher.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb';
	import { buttonVariants } from '$lib/components/ui/button';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { Label } from '$lib/components/ui/label';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { Switch } from '$lib/components/ui/switch';
	import { Toggle } from '$lib/components/ui/toggle/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import type { User } from '$lib/server/session';
	import { breadcrumbs } from '$lib/stores';

	import type { LayoutData } from './$types';

	let {
		data,
		children
	}: {
		data: LayoutData;
		children: Snippet;
	} = $props();

	// Computed values
	const current = $derived($breadcrumbs.at(-1));

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const envClient = createClient(EnvironmentService, transport);
	const resourceClient = createClient(ResourceService, transport);

	let workspaces = $state<TenantOtterscaleIoV1Alpha1Workspace[]>([]);
	let scopes = $state<Scope[]>([]);
	let previousScope = $state<string>('');
	let invalidScope = $state<string>('');
	let activeScope = $derived(page.params.scope || previousScope || 'OtterScale');

	async function fetchWorkspaces() {
		try {
			const response = await resourceClient.list({
				cluster: 'aaa',
				group: 'tenant.otterscale.io',
				version: 'v1alpha1',
				resource: 'workspaces',
				labelSelector: 'user.otterscale.io/' + data.user.sub
			});
			workspaces = response.items.map((item) => item.object as TenantOtterscaleIoV1Alpha1Workspace);
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes = response.scopes.filter((scope) => scope.name !== 'cos');
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function handleScopeOnSelect(index: number) {
		const scope = scopes[index];
		if (!scope) return;

		await goto(resolve('/(auth)/scope/[scope]', { scope: scope.name }));
	}

	async function initialize(scope: string) {
		try {
			await fetchScopes();
			// Validate scope: if not "OtterScale" and not in the scopes list, redirect to "OtterScale"
			const isValidScope = scope === 'OtterScale' || scopes.some((s) => s.name === scope);
			if (!isValidScope) {
				invalidScope = scope;
				await goto(resolve('/(auth)/scope/[scope]', { scope: 'OtterScale' }));
				return;
			}
			// Show appropriate toast based on whether we were redirected from an invalid scope
			if (invalidScope) {
				toast.warning(
					m.scope_not_found_redirect({ invalid_scope: invalidScope, scope: 'OtterScale' })
				);
				invalidScope = '';
			} else {
				toast.success(m.switch_scope({ name: scope }));
			}
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchWorkspaces();

			isMounted = true;
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	});

	$effect(() => {
		if (activeScope && activeScope !== previousScope) {
			previousScope = activeScope;
			initialize(activeScope);
		}
	});

	let next = $state(false);
</script>

<svelte:head>
	<title>{current ? `${current.title} - OtterScale` : 'OtterScale'}</title>
</svelte:head>

<Sidebar.Provider>
	<Sidebar.Root collapsible="icon" variant="inset" class="p-3">
		<Sidebar.Header>
			{#if isMounted}
				<WorkspaceSwitcher {workspaces} user={data.user} />
			{:else}
				<Sidebar.MenuButton size="lg">
					<Skeleton class="size-8 bg-foreground/10" />
					<div class="space-y-2">
						<Skeleton class="h-3 w-36 bg-foreground/10" />
						<Skeleton class="h-2 w-12 bg-foreground/10" />
					</div>
				</Sidebar.MenuButton>
			{/if}
		</Sidebar.Header>
		<Sidebar.Content class="gap-2">
			<NavDashboard items={navData.overview} />
			{#if next}
				<NavMain label="AI Studio" items={navData.aiStudio} />
				<NavMain label="Applications" items={navData.applications} />
				<NavMain label="Resources" items={navData.resources} />
				<NavMain label="Governance" items={navData.governance} />
				<NavMain label="Reliability" items={navData.reliability} />
				<NavMain label="System" items={navData.system} />
			{:else}
				<NavGeneral scope={activeScope} title={m.platform()} routes={platformRoutes(activeScope)} />
				<NavGeneral scope={activeScope} title={m.global()} routes={globalRoutes()} />
			{/if}
		</Sidebar.Content>
		<Button
			class="mx-auto w-full text-xs text-muted-foreground"
			variant="link"
			href="#"
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
				<Button variant="ghost" size="icon" class="size-7" href="/">
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
