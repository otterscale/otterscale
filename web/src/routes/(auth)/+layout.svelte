<script lang="ts">
	import { Code, ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import HouseIcon from '@lucide/svelte/icons/house';
	import type { Snippet } from 'svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { EnvironmentService, PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
	import {
		type APIResource,
		type DiscoveryRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';
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
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, premiumTier } from '$lib/stores';

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

	const tierMap = {
		[PremiumTier_Level.COMMUNITY]: m.community_tier(),
		[PremiumTier_Level.STANDARD]: m.standard_tier(),
		[PremiumTier_Level.PREMIUM]: m.premium_tier(),
		[PremiumTier_Level.ENTERPRISE]: m.enterprise_tier()
	};

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const envClient = createClient(EnvironmentService, transport);
	const resourceClient = createClient(ResourceService, transport);

	let scopes = $state<Scope[]>([]);
	let previousScope = $state<string>('');
	let invalidScope = $state<string>('');
	let activeScope = $derived(page.params.scope || previousScope || 'OtterScale');

	async function fetchScopes() {
		try {
			const response = await scopeClient.listScopes({});
			scopes = response.scopes.filter((scope) => scope.name !== 'cos');
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

	let apiResources = $state<APIResource[]>([]);
	async function fetchAPIResources() {
		try {
			const response = await resourceClient.discovery({
				cluster: activeScope
			} as DiscoveryRequest);
			apiResources = response.apiResources
				.map((resource) => ({
					...resource,
					group: resource.group ? resource.group : 'core'
				}))
				.sort((previous, next) => {
					if (previous.group !== next.group) return previous.group.localeCompare(next.group);

					if (previous.version !== next.version)
						return previous.version.localeCompare(next.version);

					return previous.kind.localeCompare(next.kind);
				});
		} catch (error) {
			console.error('Failed to fetch discoveries:', error);
		}
	}
	const mapGroupVersionKindToAPIResource = $derived(
		apiResources.reduce(
			(map, resource) => {
				map[`${resource.group}/${resource.version}/${resource.kind}`] = resource;
				return map;
			},
			{} as Record<string, APIResource>
		)
	);
	const apiResourcesByKindByVersionByGroup = $derived(
		Object.fromEntries(
			Object.entries(
				apiResources.reduce(
					(accumulation, resource) => {
						if (!accumulation[resource.group]) accumulation[resource.group] = {};

						if (!accumulation[resource.group][resource.version])
							accumulation[resource.group][resource.version] = {};

						if (!accumulation[resource.group][resource.version][resource.kind])
							accumulation[resource.group][resource.version][resource.kind] = [];

						accumulation[resource.group][resource.version][resource.kind].push(resource);
						return accumulation;
					},
					{} as Record<string, Record<string, Record<string, APIResource[]>>>
				)
			)
		)
	);

	async function handleScopeOnSelect(index: number) {
		const scope = scopes[index];
		if (!scope) return;

		await goto(resolve('/(auth)/scope/[scope]', { scope: scope.name }));
	}

	async function initialize(scope: string) {
		try {
			await fetchAPIResources();
			await fetchScopes();
			// Validate scope: if not "OtterScale" and not in the scopes list, redirect to "OtterScale"
			const isValidScope = scope === 'OtterScale' || scopes.some((s) => s.name === scope);
			if (!isValidScope) {
				invalidScope = scope;
				await goto(resolve('/(auth)/scope/[scope]', { scope: 'OtterScale' }));
				return;
			}
			await fetchEdition();
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

	$effect(() => {
		if (activeScope && activeScope !== previousScope) {
			previousScope = activeScope;
			initialize(activeScope);
		}
	});

	const enabledKinds = ['Role', 'RoleBinding', 'ClusterRole', 'ClusterRoleBinding'];
	const resources = $derived([
		{
			name: 'Manifests',
			groups: Object.entries(apiResourcesByKindByVersionByGroup).map(
				([group, apiResourcesByKindByVersion]) => ({
					name: group,
					icon: 'ph:cube',
					route: undefined,
					versions: Object.entries(apiResourcesByKindByVersion).map(
						([version, apiResourcesByKind]) => ({
							name: version,
							kinds: Object.keys(apiResourcesByKind).map((kind) => ({
								name: kind,
								enabled: enabledKinds.includes(kind),
								resources: mapGroupVersionKindToAPIResource[`${group}/${version}/${kind}`]
							}))
						})
					)
				})
			)
		}
	]);

	export function describeResourceVersion(
		version: string
	): { label: string; identifier: string } | null {
		const match = version.match(/^v(\d+)(?:(alpha|beta)(\d+))?$/);

		if (!match) {
			return null;
		}

		const [, major, stage, stageVersion] = match;

		return {
			label:
				stage === 'alpha' ? 'Experimental' : stage === 'beta' ? 'Preview' : 'General Availability',
			identifier: `v.${major}${stage ? `.${stage}` : ''}${stage ? `.${stageVersion}` : ''}`
		};
	}
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
			{#if activeScope !== 'OtterScale'}
				{#each resources as stack, index (index)}
					<Sidebar.Group>
						<Sidebar.GroupLabel>
							{stack.name}
						</Sidebar.GroupLabel>
						{#each stack.groups as group, index (index)}
							<Sidebar.Menu>
								<Collapsible.Root>
									{#snippet child({ props })}
										<Sidebar.MenuItem {...props}>
											<Sidebar.MenuButton tooltipContent={group.name}>
												<Icon icon={group.icon} />
												<Tooltip.Provider>
													<Tooltip.Root>
														<Tooltip.Trigger>
															<h4 class="max-w-36 truncate">{group.name}</h4>
														</Tooltip.Trigger>
														<Tooltip.Content>{group.name}</Tooltip.Content>
													</Tooltip.Root>
												</Tooltip.Provider>
											</Sidebar.MenuButton>
											<Collapsible.Trigger>
												{#snippet child({ props })}
													<Sidebar.MenuAction {...props} class="data-[state=open]:rotate-90">
														<Icon icon="ph:caret-right" />
														<span class="sr-only">Toggle</span>
													</Sidebar.MenuAction>
												{/snippet}
											</Collapsible.Trigger>
											<Collapsible.Content>
												{#each group.versions as version, index (index)}
													<Sidebar.MenuSub class="text-xs text-muted-foreground">
														{@const resourceVersion = describeResourceVersion(version.name)}
														{#if resourceVersion}
															<span class="flex items-center gap-1 text-muted-foreground">
																{resourceVersion.label}
																{resourceVersion.identifier}
															</span>
														{/if}
													</Sidebar.MenuSub>
													{#each version.kinds as kind, index (index)}
														<Sidebar.MenuSub>
															<Sidebar.MenuSubItem>
																<Sidebar.MenuSubButton
																	href={resolve(
																		`/${activeScope}/${kind.name}?group=${group.name}&version=${version.name}`
																	)}
																	aria-disabled={!kind.enabled}
																>
																	<Tooltip.Provider>
																		<Tooltip.Root>
																			<Tooltip.Trigger>
																				<h4 class="max-w-36 truncate">
																					{kind.name}
																				</h4>
																			</Tooltip.Trigger>
																			<Tooltip.Content>
																				{group.name}/{version.name}/{kind.name}
																			</Tooltip.Content>
																		</Tooltip.Root>
																	</Tooltip.Provider>
																</Sidebar.MenuSubButton>
															</Sidebar.MenuSubItem>
														</Sidebar.MenuSub>
													{/each}
												{/each}
											</Collapsible.Content>
										</Sidebar.MenuItem>
									{/snippet}
								</Collapsible.Root>
							</Sidebar.Menu>
						{/each}
					</Sidebar.Group>
				{/each}
			{/if}
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
