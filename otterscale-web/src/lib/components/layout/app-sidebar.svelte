<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import type { User } from 'better-auth';
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { goto } from '$app/navigation';
	import { PremiumTier, PremiumService } from '$lib/api/premium/v1/premium_pb';
	import { Essential_Type, EssentialService } from '$lib/api/essential/v1/essential_pb';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { setupScopePath } from '$lib/path';
	import {
		activeScope,
		currentCeph,
		currentKubernetes,
		tier,
		triggerUpdateScopes
	} from '$lib/stores';
	import { bookmarks, cephPaths, kubernetesPaths, routes } from './routes';
	import NavMain from './nav-main.svelte';
	import NavPrimary from './nav-primary.svelte';
	import NavSecondary from './nav-secondary.svelte';
	import NavUser from './nav-user.svelte';
	import ScopeSwitcher from './scope-switcher.svelte';

	type Props = { user: User } & ComponentProps<typeof Sidebar.Root>;

	let { user, ref = $bindable(null), ...restProps }: Props = $props();

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const premiumClient = createClient(PremiumService, transport);
	const essentialClient = createClient(EssentialService, transport);
	const scopes = writable<Scope[]>([]);

	const tierMap = {
		[PremiumTier.BASIC]: m.basic_tier(),
		[PremiumTier.ADVANCED]: m.advanced_tier(),
		[PremiumTier.ENTERPRISE]: m.enterprise_tier()
	};

	const skeletonClasses = {
		avatar: 'bg-sidebar-primary/50 size-8 rounded-lg',
		title: 'bg-sidebar-primary/50 h-3 w-[150px]',
		subtitle: 'bg-sidebar-primary/50 h-3 w-[50px]'
	};

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);

			if (response.scopes.length > 0) {
				handleScopeOnSelect(0);
			}
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function fetchEdition() {
		try {
			const response = await premiumClient.getTier({});
			tier.set(tierMap[response.tier]);
		} catch (error) {
			const connectError = error as ConnectError;
			if (connectError.code !== Code.Unimplemented) {
				console.error('Failed to fetch tier:', connectError);
			}
		}
	}

	async function fetchEssentials(uuid: string) {
		try {
			const response = await essentialClient.listEssentials({ scopeUuid: uuid });
			const { essentials } = response;

			currentCeph.set(essentials.find((e) => e.type === Essential_Type.CEPH));
			currentKubernetes.set(essentials.find((e) => e.type === Essential_Type.KUBERNETES));
		} catch (error) {
			console.error('Failed to fetch essentials:', error);
		}
	}

	async function handleScopeOnSelect(index: number) {
		const scope = $scopes[index];
		if (!scope) return;

		activeScope.set(scope);

		await fetchEssentials(scope.uuid);
		if (!$currentCeph && !$currentKubernetes) {
			toast.info(m.scope_not_configured({ name: scope.name }), {
				duration: Number.POSITIVE_INFINITY,
				action: {
					label: m.goto(),
					onClick: () => goto(setupScopePath)
				}
			});
		} else {
			toast.success(m.switch_scope({ name: scope.name }));
		}
	}

	async function initialize() {
		try {
			await Promise.all([fetchScopes(), fetchEdition()]);
		} catch (error) {
			console.error('Failed to initialize:', error);
		}
	}

	onMount(initialize);

	$effect(() => {
		if ($triggerUpdateScopes) {
			initialize();
		}
	});
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		{#if $activeScope}
			<ScopeSwitcher
				active={$activeScope}
				scopes={$scopes}
				tier={$tier}
				onSelect={handleScopeOnSelect}
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
		<NavMain {routes} {cephPaths} {kubernetesPaths} />
		<NavPrimary {bookmarks} />
		<NavSecondary class="mt-auto" />
	</Sidebar.Content>

	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
