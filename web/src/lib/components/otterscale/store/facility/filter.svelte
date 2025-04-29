<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Command from '$lib/components/ui/command';
	import * as Select from '$lib/components/ui/select/index.js';
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
	<div class="grid gap-4">
		<Label class="text-sm">CATEGORY</Label>
		<div class="grid gap-2">
			{#each candidatedCategories as category}
				<div class="flex items-center gap-2">
					<Checkbox
						checked={selectedCategories.includes(category)}
						onCheckedChange={() => {
							updateSelectedCategories(category);
							resetActivePage();
						}}
					/>
					<p class="text-xs">{category}</p>
				</div>
			{/each}
		</div>
	</div>
	<div class="flex justify-between">
		<Label class="text-sm">VERIFIED</Label>
		<Switch bind:checked={onlyVerified} />
	</div>
</div>
