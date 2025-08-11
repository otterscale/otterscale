<script lang="ts" module>
	import { type Application_Chart } from '$lib/api/application/v1/application_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';
	import { fuzzLogosIcon } from '../utils.svelte';
	import Install from './install.svelte';
</script>

<script lang="ts">
	let { chart, children }: { chart: Application_Chart; children: Snippet } = $props();

	let isDependanciesExpand = $state(false);
	let isSourcesExpand = $state(false);
	let isMaintainersExpand = $state(false);
	let isKeywordsExpand = $state(false);
</script>

<Sheet.Root>
	<Sheet.Trigger>
		{@render children()}
	</Sheet.Trigger>
	<Sheet.Content side="right" class="min-w-[23vw] p-6">
		<Sheet.Header class="flex items-start justify-between p-0">
			<Sheet.Title class="relative w-full">
				<Avatar.Root class="h-12 w-12">
					<Avatar.Image src={chart.icon} />
					<Avatar.Fallback>
						<Icon icon={fuzzLogosIcon(chart.name, 'fluent-emoji-flat:otter')} class="size-12" />
					</Avatar.Fallback>
				</Avatar.Root>
				<span class="absolute top-0 right-0">
					{#if chart.verified}
						<Badge variant="secondary">
							<Icon icon="ph:star-fill" class="h-4 w-4 fill-yellow-400 text-yellow-400" />
							Verified
						</Badge>
					{/if}
				</span>
			</Sheet.Title>
			<Sheet.Description class="text-primary text-xl">
				{chart.name}
				<span class="flex items-center gap-2">
					<Icon icon="ph:tag" class="size-4" />
					<p class="text-muted-foreground text-sm">{chart.versions[0].chartVersion}</p>
				</span>
			</Sheet.Description>
		</Sheet.Header>

		<div class="space-y-2">
			{#if chart.home}
				<div class="flex items-center gap-2">
					<Icon icon="ph:house" />
					<p class="text-muted-foreground truncate text-xs">{chart.home}</p>
				</div>
			{/if}
			{#if chart.license}
				<div class="flex items-center gap-2">
					<Icon icon="ph:identification-badge" />
					<p class="text-muted-foreground text-xs">{chart.license}</p>
				</div>
			{/if}
			{#if chart.tags}
				<div class="flex items-center gap-2">
					<Icon icon="ph:tag" />
					<p class="text-muted-foreground text-xs">{chart.tags}</p>
				</div>
			{/if}
		</div>

		<p class="text-muted-foreground my-4 h-[15vh] overflow-auto text-sm">
			{chart.description}
		</p>

		<div class="flex h-full flex-col justify-end space-y-4">
			{#if chart.dependencies && chart.dependencies.length > 0}
				<span
					class="text-muted-foreground flex items-center justify-between gap-1 text-sm font-medium"
				>
					Dependency
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
							class={cn('size-4 transition-all', isDependanciesExpand ? 'rotate-90' : '-rotate-90')}
						/>
					</Button>
				</span>
				<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto">
					{#if !isDependanciesExpand}
						{#each chart.dependencies.slice(0, 3) as dependency}
							<Badge variant="secondary" class="flex items-center gap-1 text-xs">
								<p>{dependency.name}</p>
								{#if dependency.version}
									<p>{dependency.version}</p>
								{/if}
								{#if dependency.condition}
									<p class="text-muted-foreground">{dependency.condition}</p>
								{/if}
							</Badge>
						{/each}
						{#if chart.dependencies.length > 3}
							<Badge variant="outline" class="h-fit w-fit text-xs">
								+{chart.dependencies.length - 3}
							</Badge>
						{/if}
					{:else}
						{#each chart.dependencies as dependency}
							<Badge variant="secondary" class="flex items-center gap-1 text-xs">
								<p>{dependency.name}</p>
								{#if dependency.version}
									<p>{dependency.version}</p>
								{/if}
								{#if dependency.condition}
									<p class="text-muted-foreground">{dependency.condition}</p>
								{/if}
							</Badge>
						{/each}
					{/if}
				</div>
			{/if}

			{#if chart.sources && chart.sources.length > 0}
				<span
					class="text-muted-foreground flex items-center justify-between gap-1 text-sm font-medium"
				>
					Source
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
							class={cn('size-4 transition-all', isSourcesExpand ? 'rotate-90' : '-rotate-90')}
						/>
					</Button>
				</span>
				<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto">
					{#if !isSourcesExpand}
						{#each chart.sources.slice(0, 3) as source}
							<span class="flex items-center gap-1">
								<Button variant="ghost" class="size-5" target="_blank" href={source}>
									<Icon icon="ph:link" class="size-4" />
								</Button>
								<p class="truncate text-xs">
									{source}
								</p>
							</span>
						{/each}
						{#if chart.sources.length > 3}
							<Badge variant="outline" class="group relative h-fit w-fit text-xs">
								+{chart.sources.length - 3}
							</Badge>
						{/if}
					{:else}
						{#each chart.sources as source}
							<span class="flex items-start gap-1">
								<Button variant="ghost" class="size-5" target="_blank" href={source}>
									<Icon icon="ph:link" class="size-4" />
								</Button>
								<p class="text-xs break-all">
									{source}
								</p>
							</span>
						{/each}
					{/if}
				</div>
			{/if}

			{#if chart.maintainers && chart.maintainers.length > 0}
				<span
					class="text-muted-foreground flex items-center justify-between gap-1 text-sm font-medium"
				>
					Maintainer
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
							class={cn('size-4 transition-all', isMaintainersExpand ? 'rotate-90' : '-rotate-90')}
						/>
					</Button>
				</span>
				<div class="flex max-h-[15vh] flex-col gap-1 overflow-auto">
					{#if !isMaintainersExpand}
						{#each chart.maintainers.slice(0, 3) as maintainer}
							<Badge variant="secondary" class="text-xs">
								{maintainer.name}
							</Badge>
						{/each}
						{#if chart.maintainers.length > 3}
							<Badge variant="outline" class="h-fit w-fit text-xs">
								+{chart.maintainers.length - 3}
							</Badge>
						{/if}
					{:else}
						{#each chart.maintainers as maintainer}
							<Badge variant="secondary" class="text-xs">
								{maintainer.name}
							</Badge>
						{/each}
					{/if}
				</div>
			{/if}

			{#if chart.keywords && chart.keywords.length > 0}
				<span
					class="text-muted-foreground flex items-center justify-between gap-1 text-sm font-medium"
				>
					Keyword
					<Button
						variant="outline"
						size="icon"
						class={cn('size-6', chart.keywords.length > 3 ? 'visible' : 'hidden')}
						onclick={() => {
							isKeywordsExpand = !isKeywordsExpand;
						}}
					>
						<Icon
							icon="ph:caret-left"
							class={cn('size-4 transition-all', isKeywordsExpand ? 'rotate-90' : '-rotate-90')}
						/>
					</Button>
				</span>
				<div class="flex max-h-[15vh] flex-wrap gap-1 overflow-auto">
					{#if !isKeywordsExpand}
						{#each chart.keywords.slice(0, 3) as keyword}
							<Badge variant="secondary" class="text-xs">
								{keyword}
							</Badge>
						{/each}
						{#if chart.keywords.length > 3}
							<Badge variant="outline" class="text-xs">
								+{chart.keywords.length - 3}
							</Badge>
						{/if}
					{:else}
						{#each chart.keywords as keyword}
							<Badge variant="secondary" class="text-xs">
								{keyword}
							</Badge>
						{/each}
					{/if}
				</div>
			{/if}
		</div>

		<Sheet.Footer class="p-0">
			<Install {chart} />
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
