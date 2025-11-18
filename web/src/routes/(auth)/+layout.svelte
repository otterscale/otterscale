<script lang="ts">
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import HouseIcon from '@lucide/svelte/icons/house';
	import type { Snippet } from 'svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { EnvironmentService, PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
	import { type Scope, ScopeService } from '$lib/api/scope/v1/scope_pb';
	import CreateBookmark from '$lib/components/layout/create-bookmark.svelte';
	import NavBookmark from '$lib/components/layout/nav-bookmark.svelte';
	import NavBreadcrumb from '$lib/components/layout/nav-breadcrumb.svelte';
	import NavFooter from '$lib/components/layout/nav-footer.svelte';
	import NavGeneral from '$lib/components/layout/nav-general.svelte';
	import NavUser from '$lib/components/layout/nav-user.svelte';
	import { globalRoutes, platformRoutes } from '$lib/components/layout/routes';
	import ScopeSwitcher from '$lib/components/layout/scope-switcher.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, premiumTier } from '$lib/stores';

	import type { LayoutData } from './$types';

	const EXCLUDED_SCOPES = ['cos', 'cos-dev', 'cos-lite'];

	let {
		data,
		children
	}: {
		data: LayoutData;
		children: Snippet;
	} = $props();

	// Computed values
	const current = $derived($breadcrumbs.at(-1));

	const tierMap = {
		[PremiumTier_Level.BASIC]: m.basic_tier(),
		[PremiumTier_Level.ADVANCED]: m.advanced_tier(),
		[PremiumTier_Level.ENTERPRISE]: m.enterprise_tier()
	};

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const envClient = createClient(EnvironmentService, transport);

	let scopes = $state<Scope[]>([]);
	let activeScope = $derived(page.params.scope || 'Otterscale');

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes = response.scopes.filter((scope) => !EXCLUDED_SCOPES.includes(scope.name));
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function fetchEdition() {
		try {
			const response = await envClient.getPremiumTier({});
			premiumTier.set(response);
		} catch (error) {
			const connectError = error as ConnectError;
			if (connectError.code !== Code.Unimplemented) {
				console.error('Failed to fetch tier:', connectError);
			}
		}
	}

	async function handleScopeOnSelect(index: number) {
		const scope = scopes[index];
		if (!scope) return;

		await goto(resolve('/(auth)/scope/[scope]', { scope: scope.name }));
	}

	async function initialize(scope: string) {
		try {
			activeScope = scope;
			await Promise.all([fetchScopes(), fetchEdition()]);
			toast.success(m.switch_scope({ name: scope }));
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	}

	$effect(() => {
		if (activeScope) {
			initialize(activeScope);
		}
	});
</script>

<svelte:head>
	<title>{current ? `${current.title} - OtterScale` : 'OtterScale'}</title>
</svelte:head>

<Sidebar.Provider>
	<Sidebar.Root variant="inset" collapsible="icon" class="p-3">
		<Sidebar.Header>
			<ScopeSwitcher
				active={activeScope}
				{scopes}
				tier={tierMap[$premiumTier.level]}
				onSelect={handleScopeOnSelect}
			/>
		</Sidebar.Header>

		<Sidebar.Content>
			<NavGeneral scope={activeScope} title={m.platform()} routes={platformRoutes(activeScope)} />
			<NavGeneral scope={activeScope} title={m.global()} routes={globalRoutes()} />
			<NavBookmark />
			<NavFooter class="mt-auto" />
		</Sidebar.Content>

		<Sidebar.Footer>
			<NavUser user={data.user} />
		</Sidebar.Footer>
	</Sidebar.Root>

	<Sidebar.Inset>
		<header
			class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
		>
			<div class="flex w-full items-center justify-between gap-2 px-4">
				<Sidebar.Trigger class="-ml-1 {buttonVariants({ variant: 'ghost', size: 'icon' })}" />

				<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />

				<NavBreadcrumb />

				<CreateBookmark />

				<Button href="/" variant="ghost" size="icon">
					<HouseIcon />
				</Button>
			</div>
		</header>

		<main class="flex flex-1 flex-col px-2 md:px-4 lg:px-8">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
