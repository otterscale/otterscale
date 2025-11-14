<script lang="ts">
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import type { ComponentProps } from 'svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { EnvironmentService, PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
	import { Essential_Type, OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { type Scope, ScopeService } from '$lib/api/scope/v1/scope_pb';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import type { Path } from '$lib/path';
	import type { User } from '$lib/server';
	import { bookmarks, currentCeph, currentKubernetes, premiumTier } from '$lib/stores';

	import NavBookmark from './nav-bookmark.svelte';
	import NavFooter from './nav-footer.svelte';
	import NavGeneral from './nav-general.svelte';
	import NavUser from './nav-user.svelte';
	import { globalRoutes, platformRoutes } from './routes';
	import ScopeSwitcher from './scope-switcher.svelte';

	const EXCLUDED_SCOPES = ['cos', 'cos-dev', 'cos-lite'];

	type Props = { active: string; user: User } & ComponentProps<typeof Sidebar.Root>;

	let { active, user, ref = $bindable(null), ...restProps }: Props = $props();

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const envClient = createClient(EnvironmentService, transport);
	const orchClient = createClient(OrchestratorService, transport);
	const tierMap = {
		[PremiumTier_Level.BASIC]: m.basic_tier(),
		[PremiumTier_Level.ADVANCED]: m.advanced_tier(),
		[PremiumTier_Level.ENTERPRISE]: m.enterprise_tier()
	};
	let scopes = $state<Scope[]>([]);

	async function onBookmarkDelete(path: Path) {
		bookmarks.update((currentBookmarks) =>
			currentBookmarks.filter((bookmark) => bookmark.url !== path.url)
		);
	}

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
			active = scope;
			await Promise.all([fetchScopes(), fetchEdition(), fetchEssentials(scope)]);
			toast.success(m.switch_scope({ name: scope }));
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	}

	$effect(() => {
		if (page.params.scope) {
			initialize(page.params.scope);
		}
	});
</script>

<Sidebar.Root bind:ref variant="inset" collapsible="icon" class="p-3" {...restProps}>
	<Sidebar.Header>
		<ScopeSwitcher
			{active}
			{scopes}
			tier={tierMap[$premiumTier.level]}
			onSelect={handleScopeOnSelect}
		/>
	</Sidebar.Header>

	<Sidebar.Content>
		<NavGeneral title={m.platform()} routes={platformRoutes(active)} />
		<NavGeneral title={m.global()} routes={globalRoutes()} />
		<NavBookmark bookmarks={$bookmarks} onDelete={onBookmarkDelete} />
		<NavFooter class="mt-auto" />
	</Sidebar.Content>

	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
