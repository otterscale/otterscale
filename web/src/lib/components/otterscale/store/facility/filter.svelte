<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Command from '$lib/components/ui/command';
	import type { Facility_Charm } from '$gen/api/nexus/v1/nexus_pb';

	let {
		charms,
		selectedCategories = $bindable(),
		onlyVerified = $bindable(),
		activePage = $bindable()
	}: {
		charms: Facility_Charm[];
		selectedCategories: string[];
		onlyVerified: boolean;
		activePage: number;
	} = $props();

	function resetActivePage() {
		activePage = 1;
	}

	function updateSelectedCategories(category: string) {
		if (!selectedCategories.includes(category)) {
			selectedCategories = [...selectedCategories, category];
		} else {
			selectedCategories = selectedCategories.filter((t) => t !== category);
		}
	}
	let candidatedCategories = $derived([...new Set(charms.flatMap((c) => c.categories))].sort());
</script>

<div class="flex flex-col justify-between gap-4 p-2">
	<div>
		<Label class="text-sm font-light">Category</Label>
		<Command.Root class="bg-muted/50">
			<Command.List class="max-h-[600px] p-2">
				{#each candidatedCategories as category}
					<Command.Item
						value={category}
						class="text-xs hover:cursor-pointer"
						onSelect={() => {
							updateSelectedCategories(category);
							resetActivePage();
						}}
					>
						<div class="flex items-center gap-1">
							<Checkbox class="size-3" checked={selectedCategories.includes(category)} />
							<p class="text-xs font-light">{category}</p>
						</div>
					</Command.Item>
				{/each}
			</Command.List>
		</Command.Root>
	</div>

	<div>
		<Label class="text-sm font-light">Verified</Label>

		<div class="flex justify-end rounded-lg bg-muted/50 p-2">
			<Switch bind:checked={onlyVerified} />
		</div>
	</div>
</div>
