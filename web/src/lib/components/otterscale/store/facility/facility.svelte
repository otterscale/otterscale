<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { cn } from '$lib/utils';
	import * as Table from '$lib/components/ui/table';
	import Icon from '@iconify/svelte';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';
	import {
		Nexus,
		type Facility_Charm,
		type Facility_Charm_Artifact
	} from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let {
		selectedCharm
		// selectedChartReleases
	}: {
		selectedCharm: Facility_Charm;
		// selectedChartReleases: Application_Release[] | undefined;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const artifactsStore = writable<Facility_Charm_Artifact[]>([]);
	const artifactsLoading = writable(true);
	async function fetchArtifacts() {
		try {
			const response = await client.listCharmArtifacts({
				name: selectedCharm.name
			});
			artifactsStore.set(response.artifacts);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			artifactsLoading.set(false);
		}
	}

	const charmStore = writable<Facility_Charm>();
	const charmLoading = writable(true);
	async function fetchCharm() {
		try {
			const response = await client.getCharm({
				name: selectedCharm.name
			});
			charmStore.set(response);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			charmLoading.set(false);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchArtifacts();
			await fetchCharm();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<main class="grid gap-2">
		<div class="rounded-lg bg-muted/50 p-4">
			<span class="flex items-start gap-2">
				<Avatar.Root class="h-10 w-10">
					<Avatar.Image src={$charmStore.icon} />
					<Avatar.Fallback>
						<Skeleton class="size-8 rounded-full" />
					</Avatar.Fallback>
				</Avatar.Root>
				<div class="flex flex-col items-start">
					<h1 class="text-lg">{$charmStore.name}</h1>
					<h1 class="text-sm text-muted-foreground">{$charmStore.title}</h1>
				</div>
			</span>
			<p class="justify-around hyphens-auto p-4 text-xs font-light">
				{$charmStore.description}
			</p>
		</div>

		<div
			class={cn(
				'grid max-h-[calc(70vh_-_theme(spacing.16))] gap-2 overflow-auto p-2',
				cn('[&>fieldset]:p-2'),
				cn(
					'[&>fieldset>legend]:w-full [&>fieldset>legend]:text-sm [&>fieldset>legend]:font-extralight'
				),
				cn('[&>fieldset>div]:p-2')
			)}
		>
			{#if $charmStore.website}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:house" />
						WEBSITE
					</legend>
					<div>
						<a href={$charmStore.website} target="_blank">
							<span class="flex items-center gap-1">
								<p class="break-all text-xs">{$charmStore.website}</p>
								<Icon icon="ph:arrow-square-out" />
							</span>
						</a>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.storeUrl}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:cloud" />
						STORE
					</legend>
					<div>
						<a href={$charmStore.storeUrl} target="_blank">
							<span class="flex items-center gap-1">
								<p class="break-all text-xs">{$charmStore.storeUrl}</p>
								<Icon icon="ph:arrow-square-out" />
							</span>
						</a>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.deployableOn && $charmStore.deployableOn.length > 0}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:stack" />
						DEPLOYABILITY
					</legend>
					<div>
						<span class="flex flex-wrap gap-1">
							{#each $charmStore.deployableOn as deployability}
								<Badge variant="secondary" class="w-fit text-[13px]">{deployability}</Badge>
							{/each}
						</span>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.type}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:tag" />
						TYPE
					</legend>
					<div>
						<Badge variant="outline" class="w-fit text-[13px]">{$charmStore.type}</Badge>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.categories && $charmStore.categories.length > 0}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:tag" />
						CATEGORY
					</legend>
					<div>
						<span class="flex flex-wrap gap-1">
							{#each $charmStore.categories as category}
								<Badge variant="secondary" class="w-fit text-[13px]">{category}</Badge>
							{/each}
						</span>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.license}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:identification-badge" />
						LICENSE
					</legend>
					<div>
						<Badge variant="outline" class="w-fit text-[13px]">
							{$charmStore.license}
						</Badge>
					</div>
				</fieldset>
			{/if}

			{#if $charmStore.publisher}
				<fieldset>
					<legend class="flex items-center gap-1">
						<Icon icon="ph:user" />
						PUBLISHER
					</legend>
					<div>
						<Badge variant="outline" class="w-fit text-[13px]">
							{$charmStore.publisher}
						</Badge>
					</div>
				</fieldset>
			{/if}

			{#if $artifactsStore && $artifactsStore.length > 0}
				<fieldset class="border-none">
					<legend class="flex items-center gap-1">
						<Icon icon="ph:archive" />
						RELEASE
					</legend>

					<div class="w-full overflow-x-auto">
						{@render ReadRelease($artifactsStore)}
					</div>
				</fieldset>
			{/if}
		</div>
	</main>
{:else}
	<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
		<Icon icon="ph:spinner" class="size-8 animate-spin" />
		Loading...
	</div>
{/if}

{#snippet ReadRelease(artifacts: Facility_Charm_Artifact[])}
	<Table.Root>
		<Table.Header>
			<Table.Row class="*:text-[13px] *:font-light">
				<Table.Head>CHANNEL</Table.Head>
				<Table.Head>VERSION</Table.Head>
				<Table.Head>REVISION</Table.Head>
				<Table.Head>BASES</Table.Head>
				<Table.Head>CREATE TIME</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each artifacts as artifact}
				<Table.Row class="border-none *:text-[13px]">
					<Table.Cell>{artifact.channel}</Table.Cell>
					<Table.Cell>{artifact.version}</Table.Cell>
					<Table.Cell>{artifact.revision}</Table.Cell>
					<Table.Cell>
						{artifact.bases.length}
						<!-- <div class="flex flex-col gap-1">
							{#each artifact.bases as base}
								<Badge variant="outline" class="w-fit"
									>{base.name} {base.channel} {base.architecture}</Badge
								>
							{/each}
						</div> -->
					</Table.Cell>
					<Table.Cell
						>{artifact.createdAt
							? formatTimeAgo(new Date(Number(artifact.createdAt.seconds) * 1000))
							: ''}</Table.Cell
					>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
{/snippet}
