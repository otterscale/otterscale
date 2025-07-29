<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import type { User } from 'better-auth';
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { PremiumEdition, PremiumService } from '$lib/api/premium/v1/premium_pb';
	import { ScopeService, type Scope } from '$lib/api/scope/v1/scope_pb';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, edition, loadingScopes, triggerUpdateScopes } from '$lib/stores';
	import { bookmarks, routes } from './routes';
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
	const scopes = writable<Scope[]>([]);

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes.set(response.scopes);

			if (response.scopes.length > 0) {
				activeScope.set(response.scopes[0]);
			}
		} catch (error) {
			console.error('Failed to fetch scopes:', error);
		}
	}

	async function fetchEdition() {
		try {
			const response = await premiumClient.getEdition({});
			edition.set(editionMap[response.edition]);
		} catch (error) {
			const connectError = error as ConnectError;
			if (connectError.code !== Code.Unimplemented) {
				console.error('Failed to fetch edition:', connectError);
			}
		}
	}

	async function initialize() {
		loadingScopes.set(true);
		try {
			await Promise.all([fetchScopes(), fetchEdition()]);
		} finally {
			loadingScopes.set(false);
		}
	}

	onMount(initialize);

	$effect(() => {
		if ($triggerUpdateScopes) {
			initialize();
		}
	});

	const editionMap = {
		[PremiumEdition.BASIC]: m.basic_edition(),
		[PremiumEdition.PLATINUM]: m.platinum_edition(),
		[PremiumEdition.GOLD]: m.gold_edition(),
		[PremiumEdition.ENTERPRISE]: m.enterprise_edition()
	};

	const skeletonClasses = {
		avatar: 'bg-sidebar-primary/50 size-8 rounded-lg',
		title: 'bg-sidebar-primary/50 h-3 w-[150px]',
		subtitle: 'bg-sidebar-primary/50 h-3 w-[50px]'
	};
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		{#if $loadingScopes}
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
			<ScopeSwitcher scopes={$scopes} edition={$edition} />
		{/if}
	</Sidebar.Header>

	<Sidebar.Content>
		<NavMain {routes} />
		<NavPrimary {bookmarks} />
		<NavSecondary class="mt-auto" />
	</Sidebar.Content>

	<Sidebar.Footer>
		<NavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
