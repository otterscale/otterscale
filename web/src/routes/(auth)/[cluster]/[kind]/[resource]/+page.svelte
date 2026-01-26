<script lang="ts" module>
	import { tv, type VariantProps } from 'tailwind-variants';
	export const typographyVariants = tv({
		variants: {
			variant: {
				// Headings
				h1: 'scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl',
				h2: 'scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0',
				h3: 'scroll-m-20 text-2xl font-semibold tracking-tight',
				h4: 'scroll-m-20 text-xl font-semibold tracking-tight',
				h5: 'scroll-m-20 text-lg font-semibold tracking-tight',
				h6: 'scroll-m-20 text-sm font-semibold tracking-tight',

				// Body text
				p: 'leading-7 [&:not(:first-child)]:mt-6',
				lead: 'text-xl text-muted-foreground',
				large: 'text-lg font-semibold',
				small: 'text-sm font-medium leading-none',
				muted: 'text-sm text-muted-foreground',

				// Code & Technical
				inline_code: 'rounded bg-muted px-2 py-1 font-mono text-sm',
				pre: 'overflow-x-auto rounded-lg bg-muted p-4 font-mono text-sm',

				// Quotes & Special
				blockquote: 'mt-6 border-l-2 border-muted-foreground pl-6 italic text-muted-foreground',
				caption: 'text-xs text-muted-foreground',

				// Lists
				ul: 'my-6 ml-6 list-disc [&>li]:mt-2',
				ol: 'my-6 ml-6 list-decimal [&>li]:mt-2'
			}
		},
		defaultVariants: {
			variant: 'p'
		}
	});
	export type TypographyVariant = VariantProps<typeof typographyVariants>['variant'];
</script>

<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Box from '@lucide/svelte/icons/box';
	import CircleCheck from '@lucide/svelte/icons/check-circle-2';
	import File from '@lucide/svelte/icons/file';
	import Gauge from '@lucide/svelte/icons/gauge';
	import Layers from '@lucide/svelte/icons/layers';
	import Network from '@lucide/svelte/icons/network';
	import Shield from '@lucide/svelte/icons/shield';
	import Users from '@lucide/svelte/icons/users';
	import X from '@lucide/svelte/icons/x';
	import Zap from '@lucide/svelte/icons/zap';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';

	const cluster = $derived(page.params.cluster ?? '');
	// const kind = $derived(page.params.kind ?? '');
	const resource = $derived(page.params.resource ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const namespace = $derived(page.url.searchParams.get('namespace') ?? '');

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	let object = $state<any>(undefined);
	let isMounted = $state(false);
	let isGetting = $state(false);
	let isDestroyed = false;

	async function GetResource() {
		if (isGetting || isDestroyed) return;
		isGetting = true;
		try {
			const response = await resourceClient.get({
				cluster,
				namespace,
				group,
				version,
				resource,
				name: `workspace-sample`
			});
			object = response.object;
		} catch (error) {
			console.error('Failed to get resource:', error);
		} finally {
			isGetting = false;
		}
	}

	onMount(async () => {
		await GetResource();
		isMounted = true;
	});
	onDestroy(() => {
		isDestroyed = true;
	});
</script>

{#if isMounted}
	<Field.Group class="pb-8">
		<Field.Set>
			<!-- Header -->
			<Item.Root class="w-full p-0">
				<Item.Media variant="image" class="bg-muted-foreground/50 p-2">
					<Layers size={48} />
				</Item.Media>
				<Item.Content>
					<Item.Description>
						<Badge variant="outline">{object.kind}</Badge>
						{object.apiVersion}
					</Item.Description>
					<Item.Title class={typographyVariants({ variant: 'h3' })}>
						{object.metadata?.name}
					</Item.Title>
					<Separator class="invisible" />
					<div class="grid grid-cols-6 gap-2">
						{#each [{ key: 'Creation Timestamp', value: new Date(object.metadata.creationTimestamp).toLocaleString('sv-SE') }, { key: 'Generation', value: object.metadata?.generation }, { key: 'Resource Version', value: object.metadata?.resourceVersion }] as item, index (index)}
							{#if item.value}
								<Item.Root>
									<Item.Content>
										<Item.Description>
											{item.key}
										</Item.Description>
										<Item.Title>
											{item.value}
										</Item.Title>
									</Item.Content>
								</Item.Root>
							{/if}
						{/each}
					</div>
				</Item.Content>
				<Item.Actions>
					<Button variant="ghost" size="icon-lg">
						<File />
					</Button>
				</Item.Actions>
				<Item.Footer class="flex flex-col items-start justify-start gap-2">
					<!-- Tags -->
					{@const tags = {
						Labels: object.metadata?.labels ?? {},
						Annotations: object.metadata?.annotations ?? {}
					}}
					{#each Object.entries(tags) as [key, values], index (index)}
						<Item.Root class="grid w-full grid-cols-[100px_1fr] p-0">
							<Item.Media class="relative flex w-fit items-center self-start">
								<Label>{key}</Label>
								<Badge>
									{Object.entries(values).length}
								</Badge>
							</Item.Media>
							<Item.Content class="group flex flex-row flex-wrap gap-2">
								{#each Object.entries(values) as [key, value], index (index)}
									<Badge variant="outline" class="max-w-full border">
										<p class="text-muted-foreground">{key}</p>
										<Separator orientation="vertical" class="h-1" />
										<p class="max-w-xs truncate">{value}</p>
									</Badge>
								{/each}
							</Item.Content>
						</Item.Root>
					{/each}
				</Item.Footer>
			</Item.Root>
		</Field.Set>
		<!-- Spec Section -->
		<Field.Set class="grid grid-cols-1 gap-4 ">
			<!-- Resource Quota -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
				{@const resourceQuota = object.spec?.resourceQuota?.hard ?? {}}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Gauge size={28} />
							</Item.Media>
							<Item.Content>
								<Item.Title class={typographyVariants({ variant: 'h4' })}>
									Resource Quota
								</Item.Title>
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
									Resource Quota is not configured yet. Please click the edit button at the top
									right to configure Resource Quota.
								</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{/if}
				</Card.Content>
			</Card.Root>

			<!-- Limit Range -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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
									Limit Range is not configured yet. Please click the edit button at the top right
									to configure Limit Range.
								</Empty.Description>
							</Empty.Header>
						</Empty.Root>
					{/if}
				</Card.Content>
			</Card.Root>

			<!-- Network Isolation -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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

			<!-- Users -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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
				<Card.Content class="flex flex-wrap">
					{#if users.length > 0}
						{#each users as user, index (index)}
							<Item.Root class="group w-fit hover:bg-muted/30 hover:underline">
								<Item.Content>
									<Item.Title>
										{user.name}
										<Badge variant="secondary">
											{user.role}
										</Badge>
									</Item.Title>
									<Item.Description>
										{user.subject}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/each}
					{:else}
						null
					{/if}
				</Card.Content>
			</Card.Root>
		</Field.Set>

		<Field.Set>
			<!-- Related Resources -->
			<Label class={typographyVariants({ variant: 'h4' })}>Related Resources</Label>
			{#if object.status?.namespaceRef || object.status?.authorizationPolicyRef || object.status?.peerAuthenticationRef || object.status?.roleBindingRefs?.length}
				<div class="grid gap-4 md:grid-cols-3">
					{#each [object.status?.namespaceRef, object.status?.authorizationPolicyRef, object.status?.peerAuthenticationRef, ...(object.status?.roleBindingRefs || [])] as resource (resource?.name)}
						{#if resource}
							<Item.Root class="bg-muted hover:underline">
								<Item.Media>
									<Box size={24} />
								</Item.Media>
								<Item.Content>
									<Item.Title>
										<h1>{resource.kind}</h1>
										<p class={typographyVariants({ variant: 'muted' })}>{resource.apiVersion}</p>
									</Item.Title>
									<Item.Description>
										{resource.name}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/if}
					{/each}
				</div>
			{/if}
		</Field.Set>
	</Field.Group>
{:else}{/if}
<div class="flex min-h-screen items-center justify-center bg-background px-4">
	<div class="w-full max-w-md space-y-8 text-center">
		<!-- Animated Logo/Icon -->
		<div class="flex justify-center">
			<div class="relative h-20 w-20">
				<!-- Outer rotating ring with gradient -->
				<div
					class="absolute inset-0 animate-spin rounded-full border-2 border-transparent border-t-primary border-r-primary"
					style="border-right-color: oklch(0.7 0.18 170);"
				/>
				<!-- Middle pulsing ring -->
				<div class="absolute inset-2 animate-pulse rounded-full border border-primary/40" />
				<!-- Inner accelerated spin -->
				<div
					class="absolute inset-4 animate-spin rounded-full border border-transparent border-b-primary/60"
					style="animation-direction: reverse; animation-duration: 1s;"
				/>
				<!-- Center icon -->
				<div class="absolute inset-0 flex items-center justify-center">
					<Layers class="h-10 w-10 text-primary" />
				</div>
			</div>
		</div>
		<!-- Main Text -->
		<div class="space-y-3">
			<h1 class="text-3xl font-bold tracking-tight text-foreground">Loading</h1>
			<p class="text-sm leading-relaxed text-muted-foreground">
				Preparing your workspace configuration
			</p>
		</div>
		<!-- Animated Dots -->
		<div class="flex items-center justify-center gap-1.5 py-4">
			{#each [0, 1, 2] as i}
				<div
					class="h-2 w-2 animate-bounce rounded-full bg-primary"
					style="animation-delay: {i * 150}ms; animation-duration: 1.2s;"
				/>
			{/each}
		</div>
		<!-- Progress Steps -->
		<div class="space-y-3 py-4">
			{#each [{ label: 'Fetching metadata', number: '01' }, { label: 'Loading configuration', number: '02' }, { label: 'Initializing resources', number: '03' }] as step, index}
				<div
					class="opacity-transition flex items-center gap-3"
					style="animation: fadeIn 0.6s ease-out {index * 200}ms both;"
				>
					<span class="min-w-6 text-xs font-bold tracking-wider text-primary/60">
						{step.number}
					</span>
					<div class="h-px flex-1 bg-border" />
					<span class="text-xs font-medium text-muted-foreground">
						{step.label}
					</span>
				</div>
			{/each}
		</div>
		<!-- Bottom Hint -->
		<div class="border-t border-border/30 pt-4">
			<p class="text-xs text-muted-foreground/70">✓ Your workspace is being prepared</p>
		</div>
	</div>
	<style>
		@keyframes fadeIn {
			from {
				opacity: 0;
				transform: translateX(-8px);
			}
			to {
				opacity: 1;
				transform: translateX(0);
			}
		}
		.opacity-transition {
			opacity: 1;
		}
	</style>
</div>
