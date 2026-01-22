<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Box from '@lucide/svelte/icons/box';
	import CheckCircle2 from '@lucide/svelte/icons/check-circle-2';
	import Clock from '@lucide/svelte/icons/clock';
	import Gauge from '@lucide/svelte/icons/gauge';
	import Dot from '@lucide/svelte/icons/dot';
	import File from '@lucide/svelte/icons/file';
	import Layers from '@lucide/svelte/icons/layers';
	import Network from '@lucide/svelte/icons/network';
	import Plus from '@lucide/svelte/icons/plus';
	import Zap from '@lucide/svelte/icons/zap';
	import Shield from '@lucide/svelte/icons/shield';
	import Tag from '@lucide/svelte/icons/tag';
	import Users from '@lucide/svelte/icons/users';
	import X from '@lucide/svelte/icons/x';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { Switch } from '$lib/components/ui/switch';
	import * as Table from '$lib/components/ui/table/index.js';
	import CircleCheck from '@lucide/svelte/icons/check-circle-2';

	const cluster = $derived(page.params.cluster ?? '');
	const kind = $derived(page.params.kind ?? '');
	const resource = $derived(page.params.resource ?? '');
	const group = $derived(page.url.searchParams.get('group') ?? '');
	const version = $derived(page.url.searchParams.get('version') ?? '');
	const namespace = $derived(page.url.searchParams.get('namespace') ?? '');

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	let object = $state<any>(undefined);
	let isMounted = $state(false);
	let isEditing = $state(false);
	let editedObject = $state<any>(undefined);
	let newNamespace = $state('');
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
				name: `user-9c0c49d6-fc63-478e-86a4-0d1a907c202c`
			});
			object = response.object;
			editedObject = structuredClone(response.object);
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

	function handleEdit() {
		editedObject = structuredClone(object);
		isEditing = true;
	}
	function handleCancel() {
		editedObject = structuredClone(object);
		isEditing = false;
		newNamespace = '';
	}
	function handleSave() {
		object = structuredClone(editedObject);
		isEditing = false;
		newNamespace = '';
		// 可加上 API 更新呼叫
	}
	function handleAddNamespace() {
		if (
			newNamespace.trim() &&
			!editedObject?.spec?.networkIsolation?.allowedNamespaces?.includes(newNamespace.trim())
		) {
			editedObject.spec.networkIsolation.allowedNamespaces = [
				...(editedObject.spec.networkIsolation.allowedNamespaces || []),
				newNamespace.trim()
			];
			newNamespace = '';
		}
	}
	function handleRemoveNamespace(ns: string) {
		editedObject.spec.networkIsolation.allowedNamespaces =
			editedObject.spec.networkIsolation.allowedNamespaces.filter((n: string) => n !== ns);
	}

	const isReady = $derived(
		object?.status?.conditions?.some(
			(condition: JsonValue) => condition?.type === 'Ready' && condition?.status === 'True'
		)
	);
</script>

{#if !isMounted}
	<!-- Header -->
	<div class="flex flex-col space-y-4">
		<Item.Root class="w-full">
			<Item.Media>
				<div class="rounded-lg bg-muted p-2">
					<Layers class="size-8 text-primary" />
				</div>
			</Item.Media>
			<Item.Content>
				<Item.Title class="font-mono">
					<h4 class="semibold">
						{object.kind}
					</h4>
					<h6 class="text-muted-foreground">
						{object.apiVersion}
					</h6>
				</Item.Title>
				<Item.Description>
					<h1 class="text-2xl font-bold text-foreground">
						{object.metadata?.name}
					</h1>
				</Item.Description>
				<div class="flex flex-wrap items-center gap-1 font-mono text-xs **:text-muted-foreground">
					<span class="flex items-center gap-1">
						Creation Timestamp
						{new Date(object.metadata.creationTimestamp).toLocaleString('sv-SE')}
					</span>
					<Dot size={12} />
					<span class="flex items-center gap-1">
						Generation
						{object.metadata.generation}
					</span>
					<Dot size={12} />
					<span class="flex items-center gap-1">
						Resource Version
						{object.metadata.resourceVersion}
					</span>
				</div>
			</Item.Content>
			<Item.Actions>
				<Button variant="ghost" size="icon">
					<File />
				</Button>
			</Item.Actions>
			<Item.Footer class="justift-start flex flex-col items-start">
				<!-- Labels -->
				<div class="grid grid-cols-[auto_1fr] gap-2">
					<div class="relative">
						<Label class="self-start p-1 text-xs font-black">Labels</Label>
						<Badge
							class="absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full border-background text-[10px]"
						>
							{Object.entries(object.metadata?.labels).length}
						</Badge>
					</div>
					<div class="flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.labels).slice(0, 3) as [key, value], index (index)}
							<Badge variant="outline" class="bg-secondary/30 font-mono text-xs">
								<p class="text-muted-foreground">{key}</p>
								<Separator orientation="vertical" class="h-1" />
								<p class="text-primary">{value}</p>
							</Badge>
						{/each}
						{#if Object.entries(object.metadata?.labels).length > 3}
							+{Object.entries(object.metadata?.labels).length - 3}
						{/if}
					</div>
				</div>
				<!-- Annotations -->
				<div class="grid grid-cols-[auto_1fr] gap-2">
					<div class="relative">
						<Label class="self-start p-1 text-xs font-black">Annotations</Label>
						<Badge
							class="absolute -top-1 -right-1 flex size-4 items-center justify-center rounded-full border-background text-[10px]"
						>
							{Object.entries(object.metadata?.annotations).length}
						</Badge>
					</div>
					<div class="flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.annotations).slice(0, 3) as [key, value], index (index)}
							<Badge variant="outline" class="bg-secondary/30 font-mono text-xs">
								<p class="text-muted-foreground">{key}</p>
								<Separator orientation="vertical" class="h-1" />
								<p class="max-w-3xs truncate text-primary">{value}</p>
							</Badge>
						{/each}
						{#if Object.entries(object.metadata?.annotations).length > 3}
							+{Object.entries(object.metadata?.annotations).length - 3}
						{/if}
					</div>
				</div>
			</Item.Footer>
		</Item.Root>
	</div>
	<!-- Spec Section -->
	<div class="grid grid-cols-3 gap-4 p-4">
		<div class="flex h-full flex-col gap-4">
			<!-- Resource Quota -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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
					<Card.Action>
						<Badge>{Object.keys(object.spec.resourceQuota?.hard).length}</Badge>
					</Card.Action>
				</Card.Header>
				<Card.Content class="flex-1">
					<Table.Root>
						<Table.Body>
							{#each Object.entries(object.spec.resourceQuota?.hard) as [resource, limit]}
								<Table.Row class="border-none">
									<Table.Cell class="text-muted-foreground">{resource}</Table.Cell>
									<Table.Cell class="text-end font-mono text-primary">{limit}</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
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
							<Item.Description>Strict network isolation mode</Item.Description>
						</Item.Content>
						<Item.Actions>
							{#if object.spec.networkIsolation.enabled}
								<CircleCheck size={32} />
							{:else}
								<X size={32} class="text-destructive" />
							{/if}
						</Item.Actions>
					</Item.Root>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Limit Range -->
		<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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
				<Card.Action>
					<Badge>{Object.keys(object.spec.limitRange?.limits).length}</Badge>
				</Card.Action>
			</Card.Header>
			<Card.Content class="flex-1">
				{#each object.spec.limitRange?.limits as limit, index (index)}
					<Table.Root class="caption-top">
						<Table.Caption class="text-start">
							<h4 class="font-mono text-xs font-black">{limit.type}</h4>
						</Table.Caption>
						<Table.Body>
							{#each Object.entries(limit) as [key, value]}
								{#if key !== 'type'}
									<Table.Row class="border-none">
										<Table.Cell class="text-muted-foreground">{key}</Table.Cell>
										<Table.Cell class="text-end font-mono text-primary">
											{#if typeof value === 'string'}
												{value}
											{:else if typeof value === 'object' && value}
												{#each Object.entries(value) as [k, v]}
													<Badge variant="outline" class="mr-1 bg-secondary/30 font-mono text-xs">
														<span class="text-muted-foreground">{k}</span>
														<Separator orientation="vertical" class="mx-1 h-1" />
														<span class="text-primary">{v}</span>
													</Badge>
												{/each}
											{/if}
										</Table.Cell>
									</Table.Row>
								{/if}
							{/each}
						</Table.Body>
					</Table.Root>
					<Separator class="my-2 last:hidden" />
				{/each}
			</Card.Content>
		</Card.Root>

		<div class="flex h-full flex-col gap-4">
			<!-- Allowed Namespaces (Resource Quota style) -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
				<Card.Header>
					<Card.Title>
						<Item.Root class="p-0">
							<Item.Media>
								<Layers size={20} />
							</Item.Media>
							<Item.Content>
								<Item.Title>
									<h3 class="text-sm font-semibold text-foreground">Allowed Namespaces</h3>
								</Item.Title>
							</Item.Content>
						</Item.Root>
					</Card.Title>
					<Card.Action>
						<Badge>{object.spec.networkIsolation.allowedNamespaces.length}</Badge>
					</Card.Action>
				</Card.Header>
				<Card.Content class="flex-1">
					<div class="flex flex-wrap gap-1">
						{#each object.spec.networkIsolation.allowedNamespaces as namespace}
							<Badge variant="outline" class="text-muted-foreground">
								<Network class="size-3" />
								{namespace}
							</Badge>
						{/each}
					</div>
				</Card.Content>
			</Card.Root>
			<!-- Users -->
			<Card.Root class="flex h-full flex-col border-0 bg-muted/50 shadow-none">
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
						<Badge>{object.spec.users.length}</Badge>
					</Card.Action>
				</Card.Header>
				<Card.Content class="flex-1">
					{#each object.spec.users as user}
						<Item.Root>
							<Item.Content>
								<Item.Title>
									{user.name}
								</Item.Title>
								<Item.Description>
									{user.subject}
								</Item.Description>
							</Item.Content>
							<Item.Actions>
								<Badge variant="outline" class="font-mono text-xs">
									{user.role}
								</Badge>
							</Item.Actions>
						</Item.Root>
					{/each}
				</Card.Content>
			</Card.Root>
		</div>
	</div>

	<div class="space-y-4 p-4">
		<!-- Status Conditions -->
		<!-- <section class="rounded-xl border border-border bg-card p-6">
			<h2 class="mb-6 text-sm font-semibold tracking-wider text-muted-foreground uppercase">
				Status Conditions
			</h2>
			<div class="space-y-3">
				{#each object.status?.conditions as condition, index}
					<div
						class="flex flex-col gap-4 rounded-lg border border-accent/20 bg-accent/5 p-4 md:flex-row md:items-center"
					>
						<div class="flex items-center gap-3">
							<div class="rounded-full bg-accent/10 p-2">
								<CheckCircle2 class="h-5 w-5 text-accent" />
							</div>
							<div>
								<p class="font-semibold text-foreground">{condition.type}</p>
								<p class="text-xs text-muted-foreground">{condition.reason}</p>
							</div>
						</div>
						<div class="text-sm text-muted-foreground md:ml-auto">{condition.message}</div>
					</div>
				{/each}
			</div>
		</section> -->

		<!-- Related Resources -->
		{#if object.status?.namespaceRef || object.status?.authorizationPolicyRef || object.status?.peerAuthenticationRef || object.status?.roleBindingRefs?.length}
			<section>
				<h2 class="mb-4 text-sm font-semibold tracking-wider text-muted-foreground uppercase">
					Related Resources
				</h2>
				<div class="grid gap-4 md:grid-cols-3">
					{#each [object.status?.namespaceRef, object.status?.authorizationPolicyRef, object.status?.peerAuthenticationRef, ...(object.status?.roleBindingRefs || [])] as resource (resource?.name)}
						{#if resource}
							<Item.Root>
								<Item.Media variant="icon">
									<Box class="h-5 w-5 text-primary" />
								</Item.Media>
								<Item.Content>
									<Item.Title>
										<h1 class="text-sm font-medium">{resource.kind}</h1>
										<p class="font-mono text-xs text-muted-foreground">{resource.apiVersion}</p>
									</Item.Title>
									<Item.Description>
										{resource.name}
									</Item.Description>
								</Item.Content>
							</Item.Root>
						{/if}
					{/each}
				</div>
			</section>
		{/if}
	</div>
{:else}
	<div class="py-12 text-center text-muted-foreground">Loading...</div>
{/if}
