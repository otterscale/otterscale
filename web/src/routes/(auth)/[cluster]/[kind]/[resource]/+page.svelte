<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Box from '@lucide/svelte/icons/box';
	import Boxes from '@lucide/svelte/icons/boxes';
	import CircleCheck from '@lucide/svelte/icons/check-circle-2';
	import Cylinder from '@lucide/svelte/icons/cylinder';
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
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import * as Table from '$lib/components/ui/table/index.js';

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
			console.log(object);
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
	<div class="space-y-8 pb-8">
		<!-- Header -->
		<Item.Root class="w-full space-y-2 p-0">
			<Item.Media>
				<div class="rounded-lg bg-muted p-2">
					<Layers class="size-8 text-primary" />
				</div>
			</Item.Media>
			<Item.Content>
				<Item.Title>
					<h1 class="text-2xl font-bold text-foreground">
						{object.metadata?.name}
					</h1>
				</Item.Title>
				<Item.Description>
					<div class="flex items-center gap-2">
						<h4 class="semibold">
							{object.kind}
						</h4>
						<h6 class="text-muted-foreground">
							{object.apiVersion}
						</h6>
					</div>
				</Item.Description>
				<div class="my-4 grid grid-cols-6 gap-2 text-xs **:text-muted-foreground">
					{#each [{ key: 'Creation Timestamp', value: new Date(object.metadata.creationTimestamp).toLocaleString('sv-SE') }, { key: 'Generation', value: object.metadata?.generation }, { key: 'Resource Version', value: object.metadata?.resourceVersion }] as item, index (index)}
						<Item.Root class="p-0">
							<Item.Content>
								<Item.Title class="text-xs font-semibold">{item.key}</Item.Title>
								<Item.Description class="text-xs">
									{item.value}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/each}
				</div>
			</Item.Content>
			<Item.Actions>
				<Button variant="ghost" size="icon">
					<File />
				</Button>
			</Item.Actions>
			<Item.Footer class="justift-start flex flex-col items-start">
				<div class="grid grid-cols-[auto_1fr] gap-2">
					<!-- Labels -->
					<div class="relative row-start-1 w-fit">
						<Label class="self-start p-1 text-xs font-semibold">Labels</Label>
						<Badge
							class="absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full border-background text-[10px]"
						>
							{Object.entries(object.metadata?.labels).length}
						</Badge>
					</div>
					<div class="row-start-1 flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.labels) as [key, value], index (index)}
							<Badge variant="outline" class="text-xs">
								<p class="text-muted-foreground">{key}</p>
								<Separator orientation="vertical" class="h-1" />
								<p class="text-primary">{value}</p>
							</Badge>
						{/each}
					</div>

					<!-- Annotations -->
					<div class="relative row-start-2 w-fit">
						<Label class="self-start p-1 text-xs font-semibold">Annotations</Label>
						<Badge
							class="absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full border-background text-[10px]"
						>
							{Object.entries(object.metadata?.annotations).length}
						</Badge>
					</div>
					<div class="row-start-2 flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.annotations) as [key, value], index (index)}
							<Badge variant="outline" class="text-xs">
								<p class="text-muted-foreground">{key}</p>
								<Separator orientation="vertical" class="h-1" />
								<p class="max-w-3xs truncate text-primary">{value}</p>
							</Badge>
						{/each}
					</div>
				</div>
			</Item.Footer>
		</Item.Root>
		<!-- Spec Section -->
		<div class="grid grid-cols-3 gap-4">
			<div class="flex h-full flex-col gap-4">
				<!-- Resource Quota -->
				<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
					{@const resourceQuota = object.spec?.resourceQuota?.hard ?? {}}
					<Card.Header>
						<Card.Title>
							<Item.Root class="p-0">
								<Item.Media>
									<Gauge size={20} />
								</Item.Media>
								<Item.Content>
									<Item.Title>
										<h3 class="text-sm font-semibold text-foreground">Resource Quota</h3>
									</Item.Title>
								</Item.Content>
							</Item.Root>
						</Card.Title>
					</Card.Header>
					<Card.Content class="flex-1">
						{#if Object.keys(resourceQuota).length > 0}
							<Table.Root>
								<Table.Body>
									{#each Object.entries(resourceQuota) as [resource, limit], index (index)}
										<Table.Row class="border-none">
											<Table.Cell class="text-muted-foreground">{resource}</Table.Cell>
											<Table.Cell class="text-end text-primary">{limit}</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						{:else}
							null
						{/if}
					</Card.Content>
				</Card.Root>

				<!-- Network Isolation -->
				<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
					<Card.Header>
						<Card.Title>
							<Item.Root class="p-0">
								<Item.Media>
									<Shield size={20} />
								</Item.Media>
								<Item.Content>
									<Item.Title>
										<h3 class="text-sm font-semibold text-foreground">Network Isolation</h3>
									</Item.Title>
								</Item.Content>
							</Item.Root>
						</Card.Title>
					</Card.Header>
					<Card.Content class="flex-1">
						<Item.Root class="flex w-full items-center justify-between p-0">
							<Item.Content>
								<Item.Title>Enabled</Item.Title>
								<Item.Description>Ingress traffic rules for the workspace</Item.Description>
							</Item.Content>
							<Item.Actions>
								{@const enabled = object.spec?.networkIsolation?.enabled ?? null}
								{#if enabled === true}
									<CircleCheck size={32} />
								{:else if enabled === false}
									<X size={32} class="text-destructive" />
								{:else}
									null
								{/if}
							</Item.Actions>
						</Item.Root>
					</Card.Content>
					{@const allowedNamespaces = object.spec?.networkIsolation?.allowedNamespaces ?? []}
					<Card.Header>
						<Card.Title>
							<Item.Root class="p-0">
								<Item.Content>
									<Item.Title>
										<h3 class="text-sm font-semibold text-foreground">Allowed Namespaces</h3>
									</Item.Title>
								</Item.Content>
							</Item.Root>
						</Card.Title>
						<Card.Action>
							<Badge>{allowedNamespaces.length}</Badge>
						</Card.Action>
					</Card.Header>
					<Card.Content class="flex-1">
						{#if allowedNamespaces.length > 0}
							<div class="flex flex-wrap gap-1">
								{#each allowedNamespaces as allowedNamespace, index (index)}
									<Badge variant="outline" class="text-muted-foreground">
										<Network class="size-3" />
										{allowedNamespace}
									</Badge>
								{/each}
							</div>
						{:else}
							null
						{/if}
					</Card.Content>
				</Card.Root>
			</div>

			<!-- Limit Range -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
				{@const limits = object.spec?.limitRange?.limits ?? []}
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Zap size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title>
									<h3 class="text-sm font-semibold text-foreground">Limit Range</h3>
								</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
				</Card.Header>
				<Card.Content class="-mt-3 flex-1">
					{#if limits.length > 0}
						{#each limits as limit, index (index)}
							{@const { type, ...thresholds } = limit}
							<Table.Root class="caption-top">
								<Table.Caption>
									<Item.Root class="gap-2 p-1">
										<Item.Media>
											{#if type === 'Container'}
												<Box size={20} />
											{:else if type === 'Pod'}
												<Boxes size={20} />
											{:else if type === 'PersistentVolumeClaim'}
												<Cylinder size={20} />
											{/if}
										</Item.Media>
										<Item.Content>
											<Item.Title>
												<h4 class="text-sm">{type}</h4>
											</Item.Title>
										</Item.Content>
									</Item.Root>
								</Table.Caption>
								<Table.Body>
									{#each Object.entries(thresholds) as [key, values], index (index)}
										<Table.Row class="border-none">
											<Table.Cell class="text-muted-foreground">{key}</Table.Cell>
											<Table.Cell class="text-end text-primary">
												{#if values && typeof values === 'object'}
													{#each Object.entries(values) as [key, value], index (index)}
														<Badge variant="outline" class="mr-1 bg-secondary/30 text-xs">
															<span class="text-muted-foreground">{key}</span>
															<Separator orientation="vertical" class="mx-1 h-1" />
															<span class="text-primary">{value}</span>
														</Badge>
													{/each}
												{/if}
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
							<Separator class="my-2 last:hidden" />
						{/each}
					{:else}
						null
					{/if}
				</Card.Content>
			</Card.Root>

			<div class="flex h-full flex-col gap-4">
				<!-- Users -->
				<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
					{@const users = object.spec?.users ?? []}
					<Card.Header>
						<Card.Title>
							<Item.Root class="p-0">
								<Item.Media>
									<Users size={20} />
								</Item.Media>
								<Item.Content>
									<Item.Title>
										<h3 class="text-sm font-semibold text-foreground">Users</h3>
									</Item.Title>
								</Item.Content>
							</Item.Root>
						</Card.Title>
						<Card.Action>
							<Badge>{users.length}</Badge>
						</Card.Action>
					</Card.Header>
					<Card.Content class="flex flex-wrap gap-2">
						{#if users.length > 0}
							{#each users as user, index (index)}
								<Item.Root class="w-fit hover:underline">
									<Item.Content>
										<Item.Title>
											{user.name}
											<Badge variant="outline" class="text-xs">
												{user.role}
											</Badge>
										</Item.Title>
										<Item.Description class="text-xs">
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
			</div>
		</div>

		<!-- Related Resources -->
		<Label class="text-xl font-bold">Related Resources</Label>
		{#if object.status?.namespaceRef || object.status?.authorizationPolicyRef || object.status?.peerAuthenticationRef || object.status?.roleBindingRefs?.length}
			<div class="grid gap-4 md:grid-cols-3">
				{#each [object.status?.namespaceRef, object.status?.authorizationPolicyRef, object.status?.peerAuthenticationRef, ...(object.status?.roleBindingRefs || [])] as resource (resource?.name)}
					{#if resource}
						<Item.Root class="bg-muted/50 hover:underline">
							<Item.Media>
								<Box class="size-6 text-primary" />
							</Item.Media>
							<Item.Content>
								<Item.Title>
									<h1 class="text-sm font-semibold">{resource.kind}</h1>
									<p class="text-xs text-muted-foreground">{resource.apiVersion}</p>
								</Item.Title>
								<Item.Description class="text-xs">
									{resource.name}
								</Item.Description>
							</Item.Content>
						</Item.Root>
					{/if}
				{/each}
			</div>
		{/if}
	</div>
{:else}
	<div class="py-12 text-center text-muted-foreground">Loading...</div>
{/if}
