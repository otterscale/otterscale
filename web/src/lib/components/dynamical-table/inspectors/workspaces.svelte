<script lang="ts">
	import { Box, CircleCheck, Gauge, Network, Shield, Users, X, Zap } from '@lucide/svelte';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';

	import { typographyVariants } from '$lib/components/typography/index.ts';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';

	let { object }: { object: TenantOtterscaleIoV1Alpha1Workspace } = $props();
</script>

<Field.Group class="pb-8">
	<!-- Spec Section -->
	<Field.Set class="grid grid-cols-1 gap-0 rounded-lg bg-muted/50">
		<!-- Resource Quota -->
		<Card.Root class="flex h-full flex-col border-0 bg-transparent shadow-none">
			{@const resourceQuota = object.spec?.resourceQuota?.hard ?? {}}
			<Card.Header>
				<Card.Title>
					<Item.Root class="p-0">
						<Item.Media>
							<Gauge size={28} />
						</Item.Media>
						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h4' })}>Resource Quota</Item.Title>
						</Item.Content>
					</Item.Root>
				</Card.Title>
			</Card.Header>
			<Card.Content>
				{#if Object.keys(resourceQuota).length > 0}
					<div class="grid grid-cols-1 gap-4 lg:grid-cols-3 xl:grid-cols-5">
						{#each Object.entries(resourceQuota) as [key, value], index (index)}
							<Item.Root class="w-fit p-0">
								<Item.Content class="flex gap-2">
									<Item.Description>
										{key}
									</Item.Description>
									<Item.Title class={cn(typographyVariants({ variant: 'large' }))}>
										{value}
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
								Resource Quota is not configured yet. Please click the edit button at the top right
								to configure Resource Quota.
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
							<Zap size={28} />
						</Item.Media>
						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h4' })}>Limit Range</Item.Title>
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
							<Item.Footer class="grid grid-cols-1 gap-4 lg:grid-cols-3 xl:grid-cols-5">
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
								Limit Range is not configured yet. Please click the edit button at the top right to
								configure Limit Range.
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
							<Shield size={28} />
						</Item.Media>
						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h4' })}>
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
							<Users size={28} />
						</Item.Media>
						<Item.Content>
							<Item.Title class={typographyVariants({ variant: 'h4' })}>Users</Item.Title>
						</Item.Content>
					</Item.Root>
				</Card.Title>
				<Card.Action>
					<Badge>{users.length}</Badge>
				</Card.Action>
			</Card.Header>
			<Card.Content class="flex flex-wrap gap-4">
				<!-- Users must more than one -->
				{#if users.length > 0}
					{#each users as user, index (index)}
						<Item.Root variant="muted" class="group w-fit hover:bg-muted/30 hover:underline">
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
		<div class="grid gap-4 lg:grid-cols-3">
			{#if object?.status?.authorizationPolicyRef?.name}
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
				<Item.Root variant="muted" class="hover:underline">
					<Item.Media>
						<Box size={24} />
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
					<Item.Root variant="muted" class="hover:underline">
						<Item.Media>
							<Box size={24} />
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
	</Field.Set>
</Field.Group>
