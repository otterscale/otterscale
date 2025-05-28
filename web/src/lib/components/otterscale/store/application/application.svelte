<script lang="ts">
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs/index';

	import { ReleaseDelete, ReleaseRollback, ReleaseUpdate } from '$lib/components/otterscale/index';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';

	import {
		ApplicationService,
		type Application_Chart,
		type Application_Chart_Dependency,
		type Application_Release
	} from '$gen/api/application/v1/application_pb';
	import { createClient, type Transport } from '@connectrpc/connect';

	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	import { fuzzLogosIcon } from '$lib/icon';

	const transport: Transport = getContext('transport');
	const client = createClient(ApplicationService, transport);

	const chartStore = writable<Application_Chart>();
	const chartLoading = writable(true);
	async function fetchChart() {
		try {
			const response = await client.getChart({
				name: selectedChart.name
			});
			chartStore.set(response);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			chartLoading.set(false);
		}
	}

	let {
		releases = $bindable(),
		selectedChart,
		selectedChartReleases
	}: {
		releases: Application_Release[];
		selectedChart: Application_Chart;
		selectedChartReleases: Application_Release[] | undefined;
	} = $props();

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchChart();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<main class="gap-2">
		<div class="p-4">
			<span class="flex items-center gap-2 space-x-2">
				<Avatar.Root class="h-12 w-12">
					<Avatar.Image src={$chartStore.icon} />
					<Avatar.Fallback>
						<Icon
							icon={fuzzLogosIcon($chartStore.name, 'fluent-emoji-flat:otter')}
							class="size-12"
						/>
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="space-y-1">
					<h1 class="text-lg">
						{$chartStore.name
							.split('-')
							.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
							.join(' ')}
					</h1>
					<p class="text-sm text-muted-foreground">
						{$chartStore.description}
					</p>
				</span>
			</span>
		</div>
		<Tabs.Root value="info">
			<Tabs.List>
				<Tabs.Trigger value="info">Info</Tabs.Trigger>
				<Tabs.Trigger value="release" disabled={!selectedChartReleases}>Release</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="info">
				<div
					class={cn(
						'grid max-h-[calc(77vh_-_theme(spacing.16))] gap-4 overflow-auto p-4',

						cn(
							'[&>fieldset>legend]:w-full [&>fieldset>legend]:text-sm [&>fieldset>legend]:font-extralight'
						),
						cn('[&>fieldset>div]:py-2')
					)}
				>
					{#if $chartStore.home}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:house" />
								HOME
							</legend>
							<div>
								<span class="flex items-start gap-1">
									<a href={$chartStore.home} target="_blank" class="break-all text-xs font-light">
										{$chartStore.home}
									</a>
									<Icon icon="ph:arrow-square-out" />
								</span>
							</div>
						</fieldset>
					{/if}

					{#if $chartStore.sources && $chartStore.sources.length > 0}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:cloud" />
								SOURCE
							</legend>

							<div class="grid gap-2">
								{#each $chartStore.sources as source}
									<span class="flex items-start gap-1">
										<a href={source} target="_blank" class="break-all text-xs font-light">
											{source}
										</a>
										<Icon icon="ph:arrow-square-out" />
									</span>
								{/each}
							</div>
						</fieldset>
					{/if}

					{#if $chartStore.dependencies && $chartStore.dependencies.length > 0}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:stack" />
								DEPENDENCY
							</legend>

							<div class="grid gap-2">
								{#each $chartStore.dependencies as dependency}
									<span class="flex items-start gap-1">
										<Badge variant="outline" class="w-fit text-[13px]">
											{dependency.name}: {dependency.version}
										</Badge>
										{@render ReadDependency(dependency)}
									</span>
								{/each}
							</div>
						</fieldset>
					{/if}

					{#if $chartStore.keywords && $chartStore.keywords.length > 0}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:tag" />
								KEYWORD
							</legend>

							<div>
								<span class="flex flex-wrap gap-1">
									{#each $chartStore.keywords as keyword}
										<Badge variant="secondary" class="text-[13px]">{keyword}</Badge>
									{/each}
								</span>
							</div>
						</fieldset>
					{/if}

					{#if $chartStore.tags && $chartStore.tags.length > 0}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:tag" />
								TAG
							</legend>

							<div>
								<span class="flex flex-wrap gap-1">
									{#each $chartStore.tags as tag}
										<Badge variant="secondary" class="w-fit text-[13px]">{tag}</Badge>
									{/each}
								</span>
							</div>
						</fieldset>
					{/if}

					{#if $chartStore.maintainers && $chartStore.maintainers.length > 0}
						<fieldset>
							<legend class="flex items-center gap-1">
								<Icon icon="ph:user" />
								MAINTAINER
							</legend>
							<div>
								<div class="flex flex-col gap-2">
									{#each $chartStore.maintainers as maintainer}
										<span class="flex items-start gap-1">
											<a href={maintainer.url} target="_blank">
												<Badge variant="outline" class="flex w-fit gap-2 text-[13px]">
													<p>{maintainer.name}</p>
													{#if maintainer.email}
														<p
															class="flex items-center gap-1 text-xs font-light text-muted-foreground"
														>
															<Icon icon="ph:envelope-simple" />
															{maintainer.email}
														</p>
													{/if}
												</Badge>
											</a>

											<Icon icon="ph:arrow-square-out" />
										</span>
									{/each}
								</div>
							</div>
						</fieldset>
					{/if}
				</div>
			</Tabs.Content>
			<Tabs.Content value="release">
				{#if selectedChartReleases && selectedChartReleases.length > 0}
					{@render ReadReleases(selectedChartReleases)}
				{/if}
			</Tabs.Content>
		</Tabs.Root>
	</main>
{:else}
	<ComponentLoading />
{/if}

{#snippet ReadDependency(dependency: Application_Chart_Dependency)}
	<HoverCard.Root>
		<HoverCard.Trigger>
			<Icon icon="ph:info" class="size-4 text-blue-800" />
		</HoverCard.Trigger>
		<HoverCard.Content>
			<Table.Root class="min-w-fit">
				<Table.Body
					class={cn(
						cn('[&>tr>th]:whitespace-nowrap [&>tr>th]:text-xs [&>tr>th]:font-extralight'),
						cn('[&>tr>td]:text-right [&>tr>td]:text-xs [&>tr>td]:font-light')
					)}
				>
					<Table.Row class="*:whitespace-nowrap">
						<Table.Head>NAME</Table.Head>
						<Table.Cell>{dependency.name}</Table.Cell>
					</Table.Row>
					<Table.Row class="*:whitespace-nowrap">
						<Table.Head>VERSION</Table.Head>
						<Table.Cell>{dependency.version}</Table.Cell>
					</Table.Row>
					{#if dependency.condition}
						<Table.Row class="*:whitespace-nowrap">
							<Table.Head>CONDITION</Table.Head>
							<Table.Cell>{dependency.condition}</Table.Cell>
						</Table.Row>
					{/if}
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}

{#snippet ReadReleases(selectedChartReleases: Application_Release[])}
	<div class="w-full overflow-x-auto">
		<Table.Root>
			<Table.Header>
				<Table.Row class="*:text-[13px] *:font-thin">
					<Table.Head>NAME</Table.Head>
					<Table.Head>SCOPE</Table.Head>
					<Table.Head>FACILITY</Table.Head>
					<Table.Head>NAMESPACE</Table.Head>
					<Table.Head>REVISION</Table.Head>
					<Table.Head>APPLICATION</Table.Head>
					<Table.Head>CHART</Table.Head>
					<Table.Head></Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each selectedChartReleases as release}
					<Table.Row class="border-none *:text-[13px]">
						<Table.Cell>{release.name}</Table.Cell>
						<Table.Cell>{release.scopeName}</Table.Cell>
						<Table.Cell><p class=" whitespace-nowrap">{release.facilityName}</p></Table.Cell>
						<Table.Cell>{release.namespace}</Table.Cell>
						<Table.Cell class="text-right">{release.revision}</Table.Cell>
						<Table.Cell>
							{#if release.version}
								<Badge variant="outline" class="w-fit">{release.version.applicationVersion}</Badge>
							{/if}
						</Table.Cell>
						<Table.Cell>
							{#if release.version}
								<Badge variant="outline" class="w-fit">{release.version.chartVersion}</Badge>
							{/if}
						</Table.Cell>
						<Table.Cell>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									<Button variant="ghost">
										<Icon icon="ph:dots-three-vertical" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<span class="flex items-center gap-1">
											<Icon icon="ph:arrow-clockwise" />
											<ReleaseUpdate bind:releases {release} valuesYaml={''} />
										</span>
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<span class="flex items-center gap-1">
											<Icon icon="ph:arrow-counter-clockwise" />
											<ReleaseRollback bind:releases {release} />
										</span>
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<span class="flex items-center gap-1">
											<Icon icon="ph:trash" />
											<ReleaseDelete bind:releases {release} />
										</span>
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<span class="flex items-center gap-1">
											<Icon icon="ph:rocket-launch" />
											<Button
												variant="ghost"
												onclick={() => goto(`/management/facility?scope=${release.scopeUuid}`)}
												>Facility
											</Button>
										</span>
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
{/snippet}
