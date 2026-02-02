<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		Box,
		CircleCheck,
		Gauge,
		Grid,
		HeartPulse,
		Network,
		Shield,
		Users,
		X,
		Zap
	} from '@lucide/svelte';
	import type { CoreV1ResourceQuota, TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import { getContext, onDestroy } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import {
		type DiscoveryRequest,
		type GetRequest,
		ResourceService
	} from '$lib/api/resource/v1/resource_pb';
	import { typographyVariants } from '$lib/components/typography/index.ts';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Toggle } from '$lib/components/ui/toggle/index.js';
	import { cn } from '$lib/utils';

	let { object }: { object: TenantOtterscaleIoV1Alpha1Workspace } = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	let getAbortController: AbortController | null = null;
	let mapper = new SvelteMap<string, { group: string; version: string; resource: string }>();
	async function GetResourceQuota(): Promise<CoreV1ResourceQuota> {
		getAbortController = new AbortController();
		try {
			const discovery = await resourceClient.discovery(
				{
					cluster: page.params.cluster
				} as DiscoveryRequest,
				{ signal: getAbortController.signal }
			);

			for (const apiResource of discovery.apiResources ?? []) {
				if (apiResource.resource.indexOf('/') >= 0) continue;

				const apiVersion = apiResource.group
					? `${apiResource.group}/${apiResource.version}`
					: apiResource.version;

				mapper.set(`${apiVersion}/${apiResource.kind}`, {
					group: apiResource.group,
					version: apiResource.version,
					resource: apiResource.resource
				});
			}

			const resourceQuotaReference = mapper.get(
				`${object.status?.resourceQuotaRef?.apiVersion}/${object.status?.resourceQuotaRef?.kind}`
			);

			const resourceQuotaResponse = await resourceClient.get(
				{
					cluster: page.params.cluster,
					namespace: object.status?.resourceQuotaRef?.namespace ?? '',
					group: resourceQuotaReference?.group ?? '',
					version: resourceQuotaReference?.version ?? '',
					resource: resourceQuotaReference?.resource ?? '',
					name: object.status?.resourceQuotaRef?.name ?? ''
				} as GetRequest,
				{ signal: getAbortController.signal }
			);

			return resourceQuotaResponse.object as CoreV1ResourceQuota;
		} catch {
			return {} as CoreV1ResourceQuota;
		} finally {
			if (getAbortController) {
				getAbortController.abort();
			}
		}
	}

	onDestroy(() => {
		if (getAbortController) {
			getAbortController.abort();
		}
	});

	function handleClick(
		cluster: string,
		namespace: string,
		group: string,
		version: string,
		kind: string,
		resource: string,
		name: string
	) {
		return goto(
			resolve(
				`/(auth)/${cluster}/${kind}/${resource}?group=${group}&version=${version}&name=${name}&namespace=${namespace}`
			)
		);
	}
	function getGridLayout(key: string) {
		if (key === 'limits.cpu') return 'md:row-start-1 2xl:row-start-1';
		if (key === 'limits.memory') return 'md:row-start-2 2xl:row-start-1';
		if (key === 'limits.otterscale.com/vgpu') return 'md:row-start-3 2xl:row-start-1';
		if (key === 'limits.otterscale.com/vgpumem') return 'md:row-start-4 2xl:row-start-1';
		if (key === 'limits.otterscale.com/vgpumem-percentage') return 'md:row-start-5 2xl:row-start-1';

		if (key === 'requests.cpu') return 'md:row-start-1 2xl:row-start-2';
		if (key === 'requests.memory') return 'md:row-start-2 2xl:row-start-2';
		if (key === 'requests.otterscale.com/vgpu') return 'md:row-start-3 2xl:row-start-2';
		if (key === 'requests.otterscale.com/vgpumem') return 'md:row-start-4 2xl:row-start-2';
		if (key === 'requests.otterscale.com/vgpumem-percentage')
			return 'md:row-start-5 2xl:row-start-2';

		return '2xl:row-start-3 md:row-start-6';
	}
	let isResourceQuotasGrid = $state(false);
</script>

{#await GetResourceQuota()}
	<Field.Group class="pb-8">
		<Field.Set>
			{#each Array(13).keys() as index (index)}
				{#if index % 2 === 0}
					{#if index % 3 !== 0}
						{#if index % 5 === 0}
							{#if index % 7 !== 0}
								{#if index % 11 === 0}
									<Skeleton class="h-1 w-full" />
								{:else}
									<Skeleton class="h-11 w-5/6" />
								{/if}
							{:else}
								<Skeleton class="h-7 w-4/5" />
							{/if}
						{:else}
							<Skeleton class="h-5 w-3/4" />
						{/if}
					{:else}
						<Skeleton class="h-3 w-2/3" />
					{/if}
				{:else}
					<Skeleton class="h-2 w-1/2" />
				{/if}
			{/each}
		</Field.Set>
	</Field.Group>
{:then resourceQuota}
	<Field.Group class="pb-8">
		<!-- Spec Section -->
		<Field.Set class="grid grid-cols-1 gap-0 rounded-lg bg-muted/50">
			<!-- Status Conditions -->
			<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
				{@const conditions = object.status?.conditions ?? []}
				{@const readyCondition = conditions.find((condition) => condition.type === 'Ready')}
				{@const isReady = readyCondition?.status === 'True' ? true : false}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<HeartPulse size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h6' })}>Status</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
				</Card.Header>
				<Card.Content>
					{#if conditions.length > 0}
						{#if isReady}
							<Item.Root class={cn('p-0', !isReady ? '**:text-destructive' : '**:text-none')}>
								<Item.Content>
									<Item.Title class={typographyVariants({ variant: 'h6' })}>
										{readyCondition?.reason}
									</Item.Title>
									<Item.Description>{readyCondition?.message}</Item.Description>
								</Item.Content>
							</Item.Root>
						{:else}
							{#each [...conditions].sort((p, n) => {
								const previous = new Date(p.lastTransitionTime ?? 0).getTime();
								const next = new Date(n.lastTransitionTime ?? 0).getTime();
								return next - previous;
							}) as condition, index (index)}
								<Item.Root
									class={cn(
										'p-0',
										condition.status === 'False' ? '**:text-destructive' : '**:text-none'
									)}
								>
									<Item.Content>
										<Item.Title class={typographyVariants({ variant: 'h6' })}>
											{condition?.reason}
										</Item.Title>
										<Item.Description class="max-w-3xs">{condition?.message}</Item.Description>
									</Item.Content>
								</Item.Root>
							{/each}
						{/if}
					{:else}
						<Empty.Root class="h-full">
							<Empty.Header>
								<Empty.Media variant="icon">
									<HeartPulse />
								</Empty.Media>
								<Empty.Title>No Status Available</Empty.Title>
								<Empty.Description>
									There are no status conditions to display for this workspace.
								</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{/if}
				</Card.Content>
			</Card.Root>
			<div class="px-4">
				<Separator class="py-px" />
			</div>
			<!-- Resource Quota -->
			<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
				{@const resourceQuotaHard = object.spec?.resourceQuota?.hard ?? {}}
				{@const resourceQuotaUsed: any = resourceQuota?.status?.used ?? {}}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Gauge size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h6' })}>Resource Quota</Item.Title
								>
							</Item.Content>
							<Item.Actions>
								<Toggle
									aria-label="Toggle Resource Quota Grid"
									size="sm"
									onclick={() => {
										isResourceQuotasGrid = !isResourceQuotasGrid;
									}}
									class="data-[state=on]:*:[svg]:fill-muted-foreground/50"
								>
									<Grid />
								</Toggle>
							</Item.Actions>
						</Item.Root>
					</Card.Title>
				</Card.Header>
				<Card.Content>
					{#if Object.keys(resourceQuotaHard).length > 0}
						<div class="grid grid-cols-1 gap-4 md:grid-cols-2 2xl:grid-cols-5">
							{#each Object.keys(resourceQuotaHard) as key, index (index)}
								<Item.Root
									class={cn('w-fit p-0', isResourceQuotasGrid ? getGridLayout(key) : 'flex')}
								>
									<Item.Content class="flex gap-2">
										<Item.Description>
											{key}
										</Item.Description>
										<Item.Title class={cn(typographyVariants({ variant: 'large' }))}>
											{resourceQuotaUsed[key] !== undefined
												? resourceQuotaUsed[key]
												: '?'}/{resourceQuotaHard[key]}
										</Item.Title>
									</Item.Content>
								</Item.Root>
							{/each}
						</div>
					{:else}
						<Empty.Root class="h-full">
							<Empty.Header>
								<Empty.Media variant="icon">
									<Gauge />
								</Empty.Media>
								<Empty.Title>No Resource Quota Configured</Empty.Title>
								<Empty.Description>
									Resource Quota is not configured yet. Please click the edit button at the top
									right to configure Resource Quota.
								</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{/if}
				</Card.Content>
			</Card.Root>
			<div class="px-4">
				<Separator class="py-px" />
			</div>
			<!-- Limit Range -->
			<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
				{@const limits = object.spec?.limitRange?.limits ?? []}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Zap size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h6' })}>Limit Range</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
				</Card.Header>
				<Card.Content class="h-full ">
					{#if limits.length > 0}
						{#each limits as limit, index (index)}
							{@const { type, ...thresholds } = limit}
							<Item.Root class="justify-between py-0 pl-0">
								<Item.Content>
									<Item.Title class="uppercase">
										{type}
									</Item.Title>
								</Item.Content>
								<Item.Footer
									class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5"
								>
									{#each Object.entries(thresholds) as [key, values], index (index)}
										{#if values && typeof values === 'object'}
											<Item.Root class="p-0">
												<Item.Content>
													<Item.Title
														class={cn('capitalize', typographyVariants({ variant: 'muted' }))}
													>
														{key}
													</Item.Title>
													<Item.Description>
														{#each Object.entries(values) as [key, value], index (index)}
															<p class="font-mono text-primary">{key}:{value}</p>
														{/each}
													</Item.Description>
												</Item.Content>
											</Item.Root>
										{/if}
									{/each}
								</Item.Footer>
							</Item.Root>
							<Separator class="my-2 last:hidden" />
						{/each}
					{:else}
						<Empty.Root class="h-full">
							<Empty.Header>
								<Empty.Media variant="icon">
									<Zap />
								</Empty.Media>
								<Empty.Title>No Limit Range Configured</Empty.Title>
								<Empty.Description>
									Limit Range is not configured yet. Please click the edit button at the top right
									to configure Limit Range.
								</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{/if}
				</Card.Content>
			</Card.Root>
			<div class="px-4">
				<Separator class="py-px" />
			</div>
			<!-- Network Isolation -->
			<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Shield size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h6' })}>
									Network Isolation
								</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
				</Card.Header>
				<Card.Content class="grid grid-cols-1 gap-4 lg:grid-cols-2">
					<Item.Root class="flex w-full items-center justify-between p-0">
						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h6' })}>Enabled</Item.Title>
							<Item.Description>
								{@const enabled = object.spec?.networkIsolation?.enabled ?? null}
								{#if enabled === true}
									<CircleCheck size={40} class="text-chart-2" />
								{:else if enabled === false}
									<X size={40} class="text-destructive" />
								{/if}
							</Item.Description>
						</Item.Content>
					</Item.Root>
					<Item.Root class="p-0">
						{@const allowedNamespaces = object.spec?.networkIsolation?.allowedNamespaces ?? []}

						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h6' })}>
								Allowed Namespaces
							</Item.Title>
							<Item.Description>
								{#if allowedNamespaces.length > 0}
									<div class="flex flex-wrap gap-1">
										{#each allowedNamespaces as allowedNamespace, index (index)}
											<Badge variant="secondary" class={typographyVariants({ variant: 'muted' })}>
												<Network class="size-3" />
												{allowedNamespace}
											</Badge>
										{/each}
									</div>
								{:else}
									<Badge variant="outline" class={typographyVariants({ variant: 'muted' })}>
										<Network class="size-3" />
										<p class="italic">No namespaces allowed</p>
									</Badge>
								{/if}
							</Item.Description>
						</Item.Content>
						<Item.Actions>
							<Badge>{allowedNamespaces.length}</Badge>
						</Item.Actions>
					</Item.Root>
				</Card.Content>
			</Card.Root>
			<div class="px-4">
				<Separator class="py-px" />
			</div>
			<!-- Users -->
			<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
				{@const users = object.spec?.users ?? []}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Users size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h6' })}>Users</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
					<Card.Action>
						<Badge>{users.length}</Badge>
					</Card.Action>
				</Card.Header>
				<Card.Content class="flex flex-wrap gap-8">
					<!-- Users must more than one -->
					{#if users.length > 0}
						{#each users as user, index (index)}
							<Item.Root class="p-0" size="sm">
								<Item.Content>
									<Item.Title>
										{user.name}
										<Badge>
											{user.role}
										</Badge>
									</Item.Title>
									<Item.Description>
										{user.subject}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/each}
					{/if}
				</Card.Content>
			</Card.Root>
		</Field.Set>

		<Field.Set>
			<!-- Related Resources -->
			<Label class={typographyVariants({ variant: 'h4' })}>Related Resources</Label>
			{#if Object.keys(object?.status ?? {}).length > 1}
				<div class="grid gap-4 xl:grid-cols-2 2xl:grid-cols-3">
					{#if object?.status?.authorizationPolicyRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.authorizationPolicyRef?.apiVersion}/${object?.status?.authorizationPolicyRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.authorizationPolicyRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.authorizationPolicyRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.authorizationPolicyRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.authorizationPolicyRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.authorizationPolicyRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.authorizationPolicyRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#if object?.status?.limitRangeRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.limitRangeRef?.apiVersion}/${object?.status?.limitRangeRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.limitRangeRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.limitRangeRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.limitRangeRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.limitRangeRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.limitRangeRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.limitRangeRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#if object?.status?.namespaceRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.namespaceRef?.apiVersion}/${object?.status?.namespaceRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.namespaceRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.namespaceRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.namespaceRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.namespaceRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.namespaceRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.namespaceRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#if object?.status?.networkPolicyRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.networkPolicyRef?.apiVersion}/${object?.status?.networkPolicyRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.networkPolicyRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.networkPolicyRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.networkPolicyRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.networkPolicyRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.networkPolicyRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.networkPolicyRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#if object?.status?.peerAuthenticationRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.peerAuthenticationRef?.apiVersion}/${object?.status?.peerAuthenticationRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.peerAuthenticationRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.peerAuthenticationRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.peerAuthenticationRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.peerAuthenticationRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.peerAuthenticationRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.peerAuthenticationRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#if object?.status?.resourceQuotaRef?.name}
						{@const reference = mapper.get(
							`${object?.status?.resourceQuotaRef?.apiVersion}/${object?.status?.resourceQuotaRef?.kind}`
						)}
						<Item.Root
							variant="muted"
							class="hover:underline"
							onclick={() => {
								handleClick(
									page.params.cluster!,
									object?.status?.resourceQuotaRef?.namespace ?? '',
									reference?.group ?? '',
									reference?.version ?? '',
									object?.status?.resourceQuotaRef?.kind ?? '',
									reference?.resource ?? '',
									object?.status?.resourceQuotaRef?.name ?? ''
								);
							}}
						>
							<Item.Media>
								<Box size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title class="flex flex-wrap">
									<h1>{object?.status?.resourceQuotaRef?.kind}</h1>
									<p class={typographyVariants({ variant: 'muted' })}>
										{object?.status?.resourceQuotaRef?.apiVersion}
									</p>
								</Item.Title>
								<Item.Description>
									{object?.status?.resourceQuotaRef?.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}

					{#each object?.status?.roleBindingRefs as roleBindingRef, index (index)}
						{#if roleBindingRef.name}
							{@const reference = mapper.get(`${roleBindingRef.apiVersion}/${roleBindingRef.kind}`)}

							<Item.Root
								variant="muted"
								class="hover:underline"
								onclick={() => {
									handleClick(
										page.params.cluster!,
										roleBindingRef.namespace ?? '',
										reference?.group ?? '',
										reference?.version ?? '',
										roleBindingRef.kind ?? '',
										reference?.resource ?? '',
										roleBindingRef.name ?? ''
									);
								}}
							>
								<Item.Media>
									<Box size={20} />
								</Item.Media>
								<Item.Content>
									<Item.Title class="flex flex-wrap">
										<h1>{roleBindingRef.kind}</h1>
										<p class={typographyVariants({ variant: 'muted' })}>
											{roleBindingRef.apiVersion}
										</p>
									</Item.Title>
									<Item.Description>
										{roleBindingRef.name}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/if}
					{/each}
				</div>
			{:else}
				<Empty.Root class="h-full">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Box />
						</Empty.Media>
						<Empty.Title>No Related Resources</Empty.Title>
						<Empty.Description>
							There are no related resources to display for this workspace.
						</Empty.Description>
					</Empty.Header>
				</Empty.Root>
			{/if}
		</Field.Set>
	</Field.Group>
{/await}
