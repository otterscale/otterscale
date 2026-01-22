<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Box from '@lucide/svelte/icons/box';
	import CheckCircle2 from '@lucide/svelte/icons/check-circle-2';
	import Clock from '@lucide/svelte/icons/clock';
	import Dot from '@lucide/svelte/icons/dot';
	import File from '@lucide/svelte/icons/file';
	import Layers from '@lucide/svelte/icons/layers';
	import Network from '@lucide/svelte/icons/network';
	import Plus from '@lucide/svelte/icons/plus';
	import Shield from '@lucide/svelte/icons/shield';
	import Tag from '@lucide/svelte/icons/tag';
	import Users from '@lucide/svelte/icons/users';
	import X from '@lucide/svelte/icons/x';
	import { getContext, onDestroy, onMount } from 'svelte';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { Switch } from '$lib/components/ui/switch';

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

{#if isMounted}
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
							class="absolute -top-1 left-full flex aspect-square h-4 min-w-4 -translate-x-2 items-center justify-center rounded-full border-background text-[10px]"
						>
							{Object.entries(object.metadata?.labels).length}
						</Badge>
					</div>
					<div class="flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.labels || {}).slice(0, 3) as [key, value], index (index)}
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
					<Label class="self-start p-1 text-xs font-black">Annotations</Label>
					<div class="flex flex-wrap gap-2">
						{#each Object.entries(object.metadata?.annotations || {}).slice(0, 3) as [key, value], index (index)}
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
	<div class="space-y-4 p-4">
		<!-- Spec Section -->
		<section
			class={`rounded-xl border ${isEditing ? 'border-primary/30 bg-primary/5' : 'border-border bg-card'} p-6 transition-colors`}
		>
			<div class="mb-6 flex items-center justify-between">
				<h2 class="text-sm font-semibold tracking-wider text-muted-foreground uppercase">Spec</h2>
				{#if isEditing}
					<Badge variant="outline" class="border-primary/30 bg-primary/10 text-primary"
						>Editing</Badge
					>
				{/if}
			</div>
			<div class="grid gap-8 md:grid-cols-2">
				<!-- Network Isolation -->
				<div>
					<div class="mb-4 flex items-center gap-2">
						<Shield class="h-5 w-5 text-primary" />
						<h3 class="font-semibold text-foreground">Network Isolation</h3>
					</div>
					<div class="space-y-4">
						<!-- Enabled Toggle -->
						<div class="flex items-center justify-between border-b border-border/50 py-3">
							<div>
								<p class="text-sm font-medium text-foreground">Enabled</p>
								<p class="mt-0.5 text-xs text-muted-foreground">Strict network isolation mode</p>
							</div>
							{#if isEditing}
								<Switch
									checked={editedObject.spec.networkIsolation.enabled}
									onCheckedChange={(e) => (editedObject.spec.networkIsolation.enabled = e.detail)}
								/>
							{:else}
								<Badge variant={object.spec.networkIsolation.enabled ? 'default' : 'secondary'}>
									{object.spec.networkIsolation.enabled ? 'Yes' : 'No'}
								</Badge>
							{/if}
						</div>
						<!-- Allowed Namespaces -->
						<div>
							<p class="mb-3 text-sm font-medium text-foreground">Allowed Namespaces</p>
							<div class="flex flex-wrap gap-2">
								{#each isEditing ? editedObject.spec.networkIsolation.allowedNamespaces : object.spec.networkIsolation.allowedNamespaces as ns}
									<Badge variant="outline" class="bg-secondary/50 px-3 py-1.5 font-mono text-xs">
										<Network class="mr-1.5 h-3 w-3 text-muted-foreground" />
										{ns}
										{#if isEditing}
											<button
												type="button"
												on:click={() => handleRemoveNamespace(ns)}
												class="ml-2 rounded-full p-0.5 transition-colors hover:bg-destructive/20 hover:text-destructive"
											>
												<X class="h-3 w-3" />
											</button>
										{/if}
									</Badge>
								{/each}
							</div>
							{#if isEditing}
								<div class="mt-3 flex gap-2">
									<Input
										placeholder="Add namespace..."
										bind:value={newNamespace}
										onkeydown={(e) => e.key === 'Enter' && handleAddNamespace()}
										class="flex-1"
									/>
									<Button
										type="button"
										variant="outline"
										size="icon"
										onclick={handleAddNamespace}
										disabled={!newNamespace.trim()}><Plus class="h-4 w-4" /></Button
									>
								</div>
							{/if}
						</div>
					</div>
				</div>
				<!-- Users -->
				<div>
					<div class="mb-4 flex items-center gap-2">
						<Users class="h-5 w-5 text-accent" />
						<h3 class="font-semibold text-foreground">Users</h3>
						<Badge variant="secondary" class="ml-auto">{object.spec.users.length}</Badge>
					</div>
					<div class="space-y-3">
						{#each isEditing ? editedObject.spec.users : object.spec.users as user}
							<div class="rounded-lg border border-border/50 bg-secondary/30 p-4">
								<div class="mb-2 flex items-center justify-between">
									<span class="font-medium text-foreground">{user.name}</span>
									<Badge variant="outline">{user.role}</Badge>
								</div>
								<p class="font-mono text-xs break-all text-muted-foreground">{user.subject}</p>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</section>
		<!-- Status Conditions -->
		<section class="rounded-xl border border-border bg-card p-6">
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
		</section>
		<!-- Related Resources (僅示意，需根據 object.status 內容調整) -->
		{#if object.status?.namespaceRef || object.status?.authorizationPolicyRef || object.status?.peerAuthenticationRef || object.status?.roleBindingRefs?.length}
			<section>
				<h2 class="mb-4 text-sm font-semibold tracking-wider text-muted-foreground uppercase">
					Related Resources
				</h2>
				<div class="grid gap-4 md:grid-cols-2">
					{#each [object.status?.namespaceRef, object.status?.authorizationPolicyRef, object.status?.peerAuthenticationRef, ...(object.status?.roleBindingRefs || [])] as resource (resource?.name)}
						{#if resource}
							<Card class="group transition-all hover:border-primary/30 hover:shadow-md">
								<CardContent class="pt-5 pb-5">
									<div class="flex items-start gap-4">
										<div
											class="rounded-lg bg-primary/10 p-2.5 transition-colors group-hover:bg-primary/15"
										>
											<Box class="h-5 w-5 text-primary" />
										</div>
										<div class="min-w-0 flex-1">
											<div class="mb-1 flex items-center gap-2">
												<Badge variant="outline" class="text-xs font-medium">{resource.kind}</Badge>
												<span class="font-mono text-xs text-muted-foreground"
													>{resource.apiVersion}</span
												>
											</div>
											<p class="text-sm font-medium break-all text-foreground">{resource.name}</p>
											{#if resource.namespace}
												<p class="mt-1 font-mono text-xs text-muted-foreground">
													ns: {resource.namespace}
												</p>
											{/if}
										</div>
									</div>
								</CardContent>
							</Card>
						{/if}
					{/each}
				</div>
			</section>
		{/if}
	</div>
{:else}
	<div class="py-12 text-center text-muted-foreground">Loading...</div>
{/if}
