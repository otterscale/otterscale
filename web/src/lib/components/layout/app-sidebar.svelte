<script lang="ts">
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import type { User } from 'better-auth';
	import type { ComponentProps } from 'svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import NavBookmark from './nav-bookmark.svelte';
	import NavFooter from './nav-footer.svelte';
	import NavGeneral from './nav-general.svelte';
	import NavUser from './nav-user.svelte';
	import ScopeSwitcher from './scope-switcher.svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { EnvironmentService, PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
	import { Essential_Type, OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths, getValidURL, type Path } from '$lib/path';
	import { applicationRoutes, platformRoutes } from '$lib/routes';
	import { activeScope, bookmarks, currentCeph, currentKubernetes, premiumTier } from '$lib/stores';

	type Props = { user: User } & ComponentProps<typeof Sidebar.Root>;

	let { user, ref = $bindable(null), ...restProps }: Props = $props();

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const envClient = createClient(EnvironmentService, transport);
	const orchClient = createClient(OrchestratorService, transport);
	const scopes = writable<Scope[]>([]);
	const trigger = writable<boolean>(false);

	const tierMap = {
		[PremiumTier_Level.BASIC]: m.basic_tier(),
		[PremiumTier_Level.ADVANCED]: m.advanced_tier(),
		[PremiumTier_Level.ENTERPRISE]: m.enterprise_tier(),
	};

	const skeletonClasses = {
		avatar: 'bg-sidebar-primary/50 size-8 rounded-lg',
		title: 'bg-sidebar-primary/50 h-3 w-[150px]',
		subtitle: 'bg-sidebar-primary/50 h-3 w-[50px]',
	};

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);
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

	async function fetchEssentials(uuid: string) {
		try {
			const response = await orchClient.listEssentials({ scopeUuid: uuid });
			const { essentials } = response;

			currentCeph.set(essentials.find((e) => e.type === Essential_Type.CEPH));
			currentKubernetes.set(essentials.find((e) => e.type === Essential_Type.KUBERNETES));
		} catch (error) {
			console.error('Failed to fetch essentials:', error);
		}
	}

	async function handleScopeOnSelect(index: number, home: boolean = false) {
		const scope = $scopes[index];
		if (!scope) return;

		// Set store and fetch essentials
		activeScope.set(scope);
		await fetchEssentials(scope.uuid);

		// Show success feedback
		toast.success(m.switch_scope({ name: scope.name }));

		// Go home
		if (home) {
			goto(dynamicPaths.scope(scope.name).url);
			return;
		}

		// Navigate to new url
		const url = getValidURL(page.url.pathname, scope.name, $currentCeph?.name, $currentKubernetes?.name);
		goto(url);
	}

	async function initialize() {
		try {
			await Promise.all([fetchScopes(), fetchEdition()]);
			const index = Math.max(
				$scopes.findIndex((scope) => scope.name == page.params.scope),
				0,
			);
			handleScopeOnSelect(index);
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	}

	async function onBookmarkDelete(path: Path) {
		bookmarks.update((currentBookmarks) => currentBookmarks.filter((bookmark) => bookmark.url !== path.url));
	}

	onMount(initialize);

	$effect(() => {
		if ($trigger) {
			initialize();
			trigger.set(false);
		}
	});
</script>

<Sidebar.Root bind:ref variant="inset" collapsible="icon" class="p-3" {...restProps}>
	<Sidebar.Header>
		{#if $activeScope}
			<ScopeSwitcher
				active={$activeScope}
				scopes={$scopes}
				tier={tierMap[$premiumTier.level]}
				onSelect={handleScopeOnSelect}
				{trigger}
			/>
		{:else}
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
		{/if}
	</Sidebar.Header>

	<Sidebar.Content>
		<NavGeneral title={m.applications()} routes={applicationRoutes(page.params.scope)} />
		<NavGeneral title={m.platform()} routes={platformRoutes(page.params.scope)} />
		<NavBookmark bookmarks={$bookmarks} onDelete={onBookmarkDelete} />
		<NavFooter class="mt-auto" />
	</Sidebar.Content>

	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
