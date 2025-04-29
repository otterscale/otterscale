<script lang="ts">
	import { cn } from '$lib/utils';
	import * as Command from '$lib/components/ui/command';

	import type { Facility_Charm } from '$gen/api/nexus/v1/nexus_pb';

	let {
		charms,
		searchTerm = $bindable(),
		activePage = $bindable()
	}: {
		charms: Facility_Charm[];
		searchTerm: string;
		activePage: number;
	} = $props();

	function resetActivePage() {
		activePage = 1;
	}

	function resetTerm() {
		searchTerm = '';
	}

	const charmNames = charms.map((c) => c.name);
</script>

<Command.Root>
	<Command.Input
		placeholder="Search"
		bind:value={searchTerm}
		oninput={() => {
			resetActivePage();
		}}
		onkeydown={(e) => {
			if (e.key === 'Escape') {
				resetTerm();
			}
		}}
	/>
	<Command.List
		class={cn(
			'fixed z-50 mt-10 h-fit w-fit overflow-y-auto rounded-lg bg-secondary p-2 opacity-90',
			searchTerm !== '' ? 'max-h-48' : 'max-h-0 p-0'
		)}
	>
		{#each charmNames as name}
			<Command.Item
				value={name}
				class="text-xs hover:cursor-pointer"
				onSelect={() => {
					searchTerm = name;
				}}
			>
				<p>{name}</p>
			</Command.Item>
		{/each}
	</Command.List>
</Command.Root>
