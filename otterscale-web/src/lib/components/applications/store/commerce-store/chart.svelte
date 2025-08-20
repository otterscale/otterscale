<script lang="ts">
	import {
		type Application_Chart,
		type Application_Release
	} from '$lib/api/application/v1/application_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';
	import type { Writable } from 'svelte/store';
	import Install from './chart-action-install-release.svelte';
	import { fuzzLogosIcon } from './utils';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Table from '$lib/components/custom/table';
	import Actions from './chart-actions.svelte';

	let {
		chart,
		chartReleases,
		charts = $bindable(),
		releases = $bindable(),
		children
	}: {
		chart: Application_Chart;
		chartReleases: Application_Release[] | undefined;
		charts: Writable<Application_Chart[]>;
		releases: Writable<Application_Release[]>;
		children: Snippet;
	} = $props();

	let isDependanciesExpand = $state(false);
	let isSourcesExpand = $state(false);
	let isMaintainersExpand = $state(false);
	let isKeywordsExpand = $state(false);
</script>

<Sheet.Root>
	<Sheet.Trigger>
		{@render children()}
	</Sheet.Trigger>
	<Sheet.Content side="right" class="min-w-[23vw]">
		<Sheet.Header class="bg-muted flex items-start justify-between p-6 pb-2">
			<Sheet.Title class="relative flex w-full items-start gap-2">
				<Avatar.Root class="h-12 w-12">
					<Avatar.Image src={chart.icon} />
					<Avatar.Fallback>
						<Icon icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')} class="size-12" />
					</Avatar.Fallback>
				</Avatar.Root>
				<span>
					<h3 class="font-semibold">{chart.name}</h3>
					<p class="text-muted-foreground flex items-center gap-1 text-sm">
						{chart.versions[0].chartVersion}
					</p>
				</span>
				<span class="absolute top-0 right-0">
					{#if chart.verified}
						<Badge variant="secondary">
							<Icon icon="ph:star-fill" class="h-4 w-4 fill-yellow-400 text-yellow-400" />
							Verified
						</Badge>
					{/if}
				</span>
			</Sheet.Title>
		</Sheet.Header>

		<Tabs.Root value="information">
			<Tabs.List class="-mt-4 ml-auto w-full rounded-none px-1">
				<Tabs.Trigger value="information">Information</Tabs.Trigger>
				<Tabs.Trigger value="release" disabled={!chartReleases}>Release</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="information" class="p-4">
				<div class="text-muted-foreground space-y-4 p-4 text-sm">
					{#if chart.license}
						<div class="flex items-center gap-2">
							<Icon icon="ph:identification-badge" />
							{chart.license}
						</div>
					{/if}

					<p class="max-h-[10vh] overflow-auto">
						{chart.description}
					</p>

					{#if chart.keywords && chart.keywords.length > 0}
						<span class="flex items-start gap-2">
							<div class="flex flex-wrap gap-1 overflow-auto">
								{chart.tags}
								{#each chart.keywords as keyword}
									<Badge variant="outline" class="text-xs">
										{keyword}
									</Badge>
								{/each}
							</div>
						</span>
					{/if}
				</div>

				<div class="text-muted-foreground space-y-4 p-4 text-sm">
					{#if chart.dependencies && chart.dependencies.length > 0}
						<div class="space-y-1">
							<span class="flex items-center justify-between gap-1">
								<span class="flex items-center gap-2">
									<Icon icon="ph:stack" />
									Dependency
								</span>

								<Button
									variant="outline"
									size="icon"
									class={cn('size-6', chart.dependencies.length > 3 ? 'visible' : 'hidden')}
									onclick={() => {
										isDependanciesExpand = !isDependanciesExpand;
									}}
								>
									<Icon
										icon="ph:caret-left"
										class={cn(
											'size-4 transition-all',
											isDependanciesExpand ? 'rotate-90' : '-rotate-90'
										)}
									/>
								</Button>
							</span>
							<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto pl-6 text-xs">
								{#if !isDependanciesExpand}
									{#each chart.dependencies.slice(0, 3) as dependency}
										<span class="flex items-center gap-1">
											{dependency.name}
											{#if dependency.version}
												{dependency.version}
											{/if}
											{#if dependency.condition}
												<p class="text-muted-foreground">{dependency.condition}</p>
											{/if}
										</span>
									{/each}
									{#if chart.dependencies.length > 3}
										<Badge variant="outline" class="h-fit w-fit">
											+{chart.dependencies.length - 3}
										</Badge>
									{/if}
								{:else}
									{#each chart.dependencies as dependency}
										<span class="flex items-center gap-1">
											{dependency.name}
											{#if dependency.version}
												{dependency.version}
											{/if}
											{#if dependency.condition}
												<p class="text-muted-foreground">{dependency.condition}</p>
											{/if}
										</span>
									{/each}
								{/if}
							</div>
						</div>
					{/if}

					{#if chart.home}
						<div class="space-y-1">
							<span class="flex items-center gap-2">
								<Icon icon="ph:house" />
								Home
							</span>
							<a
								target="_blank"
								href={chart.home}
								class="truncate pl-6 text-xs underline hover:no-underline"
							>
								{chart.home}
							</a>
						</div>
					{/if}

					{#if chart.sources && chart.sources.length > 0}
						<div class="space-y-1">
							<span class="flex items-center justify-between gap-1">
								<span class="flex items-center gap-2">
									<Icon icon="ph:link" />
									Source
								</span>

								<Button
									variant="outline"
									size="icon"
									class="size-6"
									onclick={() => {
										isSourcesExpand = !isSourcesExpand;
									}}
								>
									<Icon
										icon="ph:caret-left"
										class={cn(
											'size-4 transition-all',
											isSourcesExpand ? 'rotate-90' : '-rotate-90'
										)}
									/>
								</Button>
							</span>
							<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto pl-6 text-xs">
								{#if !isSourcesExpand}
									{#each chart.sources.slice(0, 3) as source}
										<a
											target="_blank"
											href={source}
											class="hover:text-primary underline hover:no-underline"
										>
											{source}
										</a>
									{/each}
									{#if chart.sources.length > 3}
										<Badge variant="outline" class="group relative h-fit w-fit">
											+{chart.sources.length - 3}
										</Badge>
									{/if}
								{:else}
									{#each chart.sources as source}
										<a
											target="_blank"
											href={source}
											class="hover:text-primary underline hover:no-underline"
										>
											{source}
										</a>
									{/each}
								{/if}
							</div>
						</div>
					{/if}

					{#if chart.maintainers && chart.maintainers.length > 0}
						<div class="space-y-1">
							<span class="flex items-center justify-between gap-1">
								<span class="flex items-center gap-2">
									<Icon icon="ph:user" />
									Maintainer
								</span>
								<Button
									variant="outline"
									size="icon"
									class={cn('size-6', chart.maintainers.length > 3 ? 'visible' : 'hidden')}
									onclick={() => {
										isMaintainersExpand = !isMaintainersExpand;
									}}
								>
									<Icon
										icon="ph:caret-left"
										class={cn(
											'size-4 transition-all',
											isMaintainersExpand ? 'rotate-90' : '-rotate-90'
										)}
									/>
								</Button>
							</span>
							<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto pl-6 text-xs">
								{#if !isMaintainersExpand}
									{#each chart.maintainers.slice(0, 3) as maintainer}
										{maintainer.name}
									{/each}
									{#if chart.maintainers.length > 3}
										<Badge variant="outline" class="h-fit w-fit text-xs">
											+{chart.maintainers.length - 3}
										</Badge>
									{/if}
								{:else}
									{#each chart.maintainers as maintainer}
										{maintainer.name}
									{/each}
								{/if}
							</div>
						</div>
					{/if}
				</div>
			</Tabs.Content>
			<Tabs.Content value="release" class="p-4">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>
								NAME
								<Table.SubHead>NAMESPACE</Table.SubHead>
							</Table.Head>
							<Table.Head>
								CHART
								<Table.SubHead>APPLICATION</Table.SubHead>
							</Table.Head>
							<Table.Head>REVISION</Table.Head>
							<Table.Head></Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#if chartReleases}
							{#each chartReleases as release}
								<Table.Row>
									<Table.Cell>
										{release.name}
										<Table.SubCell>{release.namespace}</Table.SubCell>
									</Table.Cell>
									<Table.Cell>
										{release.version?.chartVersion}
										<Table.SubCell>{release.version?.applicationVersion}</Table.SubCell>
									</Table.Cell>
									<Table.Cell>
										{release.revision}
									</Table.Cell>
									<Table.Cell>
										<Actions {release} bind:releases />
									</Table.Cell>
								</Table.Row>
							{/each}
						{/if}
					</Table.Body>
				</Table.Root>
			</Tabs.Content>
		</Tabs.Root>

		<Sheet.Footer class="p-4">
			<Install {chart} bind:charts />
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
