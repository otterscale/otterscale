<script lang="ts">
	import { Braces } from '@lucide/svelte';
	import File from '@lucide/svelte/icons/file';
	import Layers from '@lucide/svelte/icons/layers';
	import { type Snippet } from 'svelte';
	import Monaco from 'svelte-monaco';
	import { stringify } from 'yaml';

	import { typographyVariants } from '$lib/components/typography/index.ts';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Item from '$lib/components/ui/item';
	import Label from '$lib/components/ui/label/label.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

	let {
		cluster,
		namespace,
		group,
		version,
		kind,
		resource,
		name,
		object,
		children,
		status
	}: {
		cluster: string;
		namespace: string;
		group: string;
		version: string;
		kind: string;
		resource: string;
		name: string;
		object: any;
		children?: Snippet;
		status?: Snippet;
	} = $props();
</script>

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
				{@render status?.()}
				<Separator class="invisible" />
				<div class="grid gap-2 lg:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-6">
					{#each [{ key: 'Cluster', value: cluster }, { key: 'Namespace', value: namespace }, { key: 'Creation Timestamp', value: new Date(object.metadata.creationTimestamp).toLocaleString('sv-SE') }, { key: 'Generation', value: object.metadata?.generation }, { key: 'Resource Version', value: object.metadata?.resourceVersion }] as item, index (index)}
						{#if item.value}
							<Item.Root class="p-0">
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
				<Sheet.Root>
					<Sheet.Trigger>
						<Button variant="outline" size="icon-lg">
							<File />
						</Button>
					</Sheet.Trigger>
					<Sheet.Content
						side="right"
						class="flex h-full max-w-[62vw] min-w-[50vw] flex-col gap-0 overflow-y-auto p-4"
					>
						<Sheet.Header class="shruk-0 space-y-4">
							<Sheet.Title>{name}</Sheet.Title>
							<Sheet.Description>
								{group}/{version}/{kind}/{resource}
							</Sheet.Description>
						</Sheet.Header>
						{#if object}
							<div class="h-full p-4 pt-0">
								<Monaco
									value={stringify(object)}
									options={{
										language: 'yaml',
										padding: { top: 24 },
										automaticLayout: true,
										domReadOnly: true,
										readOnly: true
									}}
									theme="vs-dark"
								/>
							</div>
						{:else}
							<Empty.Root class="m-4 bg-muted/50">
								<Empty.Header>
									<Empty.Media variant="icon">
										<Braces size={36} />
									</Empty.Media>
									<Empty.Title>No Data</Empty.Title>
									<Empty.Description>
										No data is currently available for this resource.
										<br />
										To populate this resource, please add properties or values through the resource editor.
									</Empty.Description>
								</Empty.Header>
								<Empty.Content></Empty.Content>
							</Empty.Root>
						{/if}
					</Sheet.Content>
				</Sheet.Root>
			</Item.Actions>
			<Item.Footer class="flex flex-col items-start justify-start gap-2">
				<!-- Tags -->
				{@const tags = {
					Labels: object.metadata?.labels ?? {},
					Annotations: object.metadata?.annotations ?? {}
				}}
				{#each Object.entries(tags) as [key, values], index (index)}
					{#if Object.keys(values).length > 0}
						<Item.Root class="grid w-full grid-cols-[80px_1fr] p-0">
							<Item.Media class="relative flex w-fit items-center self-start">
								<Label>{key}</Label>
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
					{/if}
				{/each}
			</Item.Footer>
		</Item.Root>
	</Field.Set>
	{@render children?.()}
</Field.Group>
