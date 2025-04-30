<script lang="ts">
	import { cn } from '$lib/utils';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import type { Application_Chart } from '$gen/api/nexus/v1/nexus_pb';

	let {
		charts,
		selectedKeywords = $bindable(),
		activePage = $bindable()
	}: {
		charts: Application_Chart[];
		selectedKeywords: string[];
		activePage: number;
	} = $props();

	function resetActivePage() {
		activePage = 1;
	}

	let filterTerm = $state('');
	function resetTerm() {
		filterTerm = '';
	}

	function updateSelectedKeywords(keyword: string) {
		if (!selectedKeywords.includes(keyword)) {
			selectedKeywords = [...selectedKeywords, keyword];
		} else {
			selectedKeywords = selectedKeywords.filter((t) => t !== keyword);
		}
	}
	let candidatedKeywords = $derived([...new Set(charts.flatMap((c) => c.keywords))].sort());

	let groupedKeywords = $derived(
		candidatedKeywords.reduce(
			(m, keyword) => {
				const firstLetter = keyword[0].toUpperCase();
				if (!m[firstLetter]) m[firstLetter] = [];
				m[firstLetter].push(keyword);
				return m;
			},
			{} as Record<string, string[]>
		)
	);

	let sortedGroupedKeywords = $derived(Object.entries(groupedKeywords).sort());
</script>

<main class="flex flex-col justify-between bg-muted/50">
	<Command.Root class="rounded-none border">
		<Command.Input
			placeholder="Filter"
			bind:value={filterTerm}
			onkeydown={(e) => {
				if (e.key === 'Escape') {
					resetTerm();
				}
			}}
		/>
		<div class="flex justify-between">
			{@render MiniMap()}
			<Command.List class="min-h-[600px] w-full overflow-y-auto">
				{#each sortedGroupedKeywords as [letter, keywords]}
					<Command.Group dir="ltr" heading={letter} data-group={letter} class="mx-1 border-b">
						{#each keywords as keyword}
							<Command.Item
								value={keyword}
								class={cn(
									'h-6 overflow-clip whitespace-nowrap text-xs',
									selectedKeywords.includes(keyword)
										? 'pointer-events-none text-muted-foreground'
										: 'hover:cursor-pointer'
								)}
								onSelect={() => {
									updateSelectedKeywords(keyword);
									resetActivePage();
								}}
							>
								{keyword}
							</Command.Item>
						{/each}
					</Command.Group>
				{/each}
			</Command.List>
		</div>
	</Command.Root>
</main>

{#snippet MiniMap()}
	<div class="flex flex-col items-center justify-end bg-muted/50">
		{#each sortedGroupedKeywords as [letter]}
			<Button
				variant="ghost"
				class="h-fit w-4 p-0 text-[10px] font-bold"
				style="height: {100 / sortedGroupedKeywords.length}%"
				onmouseover={() => {
					document
						.querySelector(`[data-group="${letter}"]`)
						?.scrollIntoView({ behavior: 'smooth', block: 'start' });
				}}
			>
				{letter}
			</Button>
		{/each}
	</div>
{/snippet}
