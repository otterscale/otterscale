<script lang="ts">
	import 'driver.js/dist/driver.css';

	import { createClient, type Transport } from '@connectrpc/connect';
	import BotIcon from '@lucide/svelte/icons/bot';
	import BoxIcon from '@lucide/svelte/icons/box';
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import CloudBackupIcon from '@lucide/svelte/icons/cloud-backup';
	import CodeIcon from '@lucide/svelte/icons/code';
	import CombineIcon from '@lucide/svelte/icons/combine';
	import CpuIcon from '@lucide/svelte/icons/cpu';
	import DatabaseIcon from '@lucide/svelte/icons/database';
	import DumbbellIcon from '@lucide/svelte/icons/dumbbell';
	import FlagIcon from '@lucide/svelte/icons/flag';
	import GaugeIcon from '@lucide/svelte/icons/gauge';
	import HardDriveIcon from '@lucide/svelte/icons/hard-drive';
	import HelpCircleIcon from '@lucide/svelte/icons/help-circle';
	import HouseIcon from '@lucide/svelte/icons/house';
	import LayoutGridIcon from '@lucide/svelte/icons/layout-grid';
	import MapIcon from '@lucide/svelte/icons/map';
	import NetworkIcon from '@lucide/svelte/icons/network';
	import PcCaseIcon from '@lucide/svelte/icons/pc-case';
	import ScaleIcon from '@lucide/svelte/icons/scale';
	import ShieldCheckIcon from '@lucide/svelte/icons/shield-check';
	import ShipIcon from '@lucide/svelte/icons/ship';
	import TelescopeIcon from '@lucide/svelte/icons/telescope';
	import TerminalIcon from '@lucide/svelte/icons/terminal';
	import UnplugIcon from '@lucide/svelte/icons/unplug';
	import UsersIcon from '@lucide/svelte/icons/users';
	import WorkflowIcon from '@lucide/svelte/icons/workflow';
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
	import { activeNamespace, activeWorkspaceName } from '$lib/stores';

	import type { LayoutData } from './$types';

	const navData = {
		overview: [
			{
				name: m.workspace(),
				url: resolve(
					`/(auth)/${page.params.scope ?? page.params.cluster}/Workspace/workspaces?group=tenant.otterscale.io&version=v1alpha1&name=${$activeWorkspaceName}`
				),
				icon: MapIcon,
				edit: false
			},
			{
				name: m.resource_quota(),
				url: resolve(
					`/(auth)/${page.params.scope ?? page.params.cluster}/ResourceQuota?group=&version=v1&namespace=${$activeNamespace}&resource=resourcequotas`
				),
				icon: BoxIcon,
				edit: false
			}
		],
		aiStudio: [
			{
				title: 'Inference',
				url: '#',
				icon: BotIcon,
				items: [
					{
						title: 'Model',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Training',
				url: '#',
				icon: DumbbellIcon,
				items: [
					{
						title: 'Finetune Job',
						url: '#',
						disabled: true
					},
					{
						title: 'Dataset',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Notebooks',
				url: '#',
				icon: TerminalIcon,
				items: [
					{
						title: 'Jupyter',
						url: '#',
						disabled: true
					}
				]
			}
		],
		applications: [
			{
				title: 'Hub',
				url: '#',
				icon: LayoutGridIcon,
				items: [
					{
						title: 'Release',
						url: '#',
						disabled: true
					},
					{
						title: 'Chart',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Cloud IDE',
				url: '#',
				icon: CodeIcon,
				items: [
					{
						title: 'Coder',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Database',
				url: '#',
				icon: DatabaseIcon,
				items: [
					{
						title: 'Postgres',
						url: '#',
						disabled: true
					},
					{
						title: 'Redis',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Workflow',
				url: '#',
				icon: WorkflowIcon,
				items: [
					{
						title: 'Pipeline',
						url: '#',
						disabled: true
					},
					{
						title: 'Task',
						url: '#',
						disabled: true
					}
				]
			}
		],
		resources: [
			{
				title: 'Workloads',
				url: '#',
				icon: FlagIcon,
				items: [
					{
						title: 'Deployment',
						url: '#',
						disabled: true
					},
					{
						title: 'Stateful Set',
						url: '#',
						disabled: true
					},
					{
						title: 'Daemon Set',
						url: '#',
						disabled: true
					},
					{
						title: 'Cron Job',
						url: '#',
						disabled: true
					},
					{
						title: 'Job',
						url: '#',
						disabled: true
					},
					{
						title: 'Pod',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Compute',
				url: '#',
				icon: CpuIcon,
				items: [
					{
						title: 'Virtual Machine',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Network',
				url: '#',
				icon: NetworkIcon,
				items: [
					{
						title: 'VPC',
						url: '#',
						disabled: true
					},
					{
						title: 'Load Balancer',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Storage',
				url: '#',
				icon: HardDriveIcon,
				items: [
					{
						title: 'Block Pool',
						url: '#',
						disabled: true
					},
					{
						title: 'File System',
						url: '#',
						disabled: true
					}
				]
			}
		],
		governance: [
			{
				title: 'Tenant',
				url: '#',
				icon: UsersIcon,
				items: [
					{
						title: 'Workspace',
						url: '#',
						disabled: true
					},
					{
						title: 'User',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Policy',
				url: '#',
				icon: ScaleIcon,
				items: [
					{
						title: 'Policy',
						url: '#',
						disabled: true
					},
					{
						title: 'Compliance',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Metering',
				url: '#',
				icon: GaugeIcon,
				items: [
					{
						title: 'Budget',
						url: '#',
						disabled: true
					},
					{
						title: 'Usage',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Audit',
				url: '#',
				icon: ShieldCheckIcon,
				items: [
					{
						title: 'Trail',
						url: '#',
						disabled: true
					},
					{
						title: 'Log',
						url: '#',
						disabled: true
					}
				]
			}
		],
		reliability: [
			{
				title: 'Telemetry',
				url: '#',
				icon: TelescopeIcon,
				items: [
					{
						title: 'Collector',
						url: '#',
						disabled: true
					},
					{
						title: 'Rule',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Recovery',
				url: '#',
				icon: CloudBackupIcon,
				items: [
					{
						title: 'Backup',
						url: '#',
						disabled: true
					},
					{
						title: 'Restore',
						url: '#',
						disabled: true
					}
				]
			}
		],
		system: [
			{
				title: 'Fleet',
				url: '#',
				icon: ShipIcon,
				items: [
					{
						title: 'Cluster',
						url: '#',
						disabled: true
					},
					{
						title: 'Config',
						url: '#',
						disabled: true
					},
					{
						title: 'Image',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Metal',
				url: '#',
				icon: PcCaseIcon,
				items: [
					{
						title: 'Server',
						url: '#',
						disabled: true
					}
				]
			},
			{
				title: 'Tunnels',
				url: '#',
				icon: UnplugIcon,
				items: [
					{
						title: 'Server',
						url: '#',
						disabled: true
					},
					{
						title: 'Client',
						url: '#',
						disabled: true
					}
				]
			}
		]
	};

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

	let activeScope = $state(page.params.scope ?? page.params.cluster ?? ''); //TODO: remove page.params.scope after route updated
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
				<WorkspaceSwitcher
					cluster={activeScope}
					{workspaces}
					user={data.user}
					onsuccess={() => fetchWorkspaces(activeScope)}
				/>
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
					<NavGeneral title={m.platform()} routes={platformRoutes(activeScope, $activeNamespace)} />
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
