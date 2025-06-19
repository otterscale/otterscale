<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import { getStartPoint, setValue } from './utils.svelte';
	import { Tabs } from '$lib/components/custom/tabs/index';

	let {
		start = $bindable(),
		end = $bindable()
	}: {
		start: Date;
		end: Date;
	} = $props();

	function startSetter(date: Date) {
		start = date;
	}
	function endSetter(date: Date) {
		end = date;
	}
</script>

<Tabs.Root value="last">
	<Tabs.List class="-mt-1">
		<Tabs.Trigger value="last">Last</Tabs.Trigger>
		<Tabs.Trigger value="previous">Previous</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="last">
		<Command.Root class="gap-2">
			<Command.Input placeholder="Search" />
			<Command.List>
				<Command.Empty>No results found.</Command.Empty>
				<Command.Group heading="minute">
					{#each [1, 5, 15, 30] as minute}
						<Command.Item
							onclick={() => setValue(startSetter, endSetter, getStartPoint(), minute * 60 * 1000)}
						>
							<Badge variant="outline">Last</Badge>
							{minute} minute{minute > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="hour">
					{#each [1, 3, 6, 12] as hour}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint(), hour * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Last</Badge>
							{hour} hour{hour > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="day">
					{#each [1, 7, 15] as day}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint(), day * 24 * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Last</Badge>
							{day} day{day > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="month">
					{#each [30, 45, 90, 180] as days}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint(), days * 24 * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Last</Badge>
							{days} days
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Tabs.Content>
	<Tabs.Content value="previous">
		<Command.Root class="gap-2">
			<Command.Input placeholder="Search" />
			<Command.List>
				<Command.Empty>No results found.</Command.Empty>
				<Command.Group heading="minute">
					{#each [1, 5, 15, 30] as minute}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint('minute'), minute * 60 * 1000)}
						>
							<Badge variant="outline">Previous</Badge>
							{minute} minute{minute > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="hour">
					{#each [1, 3, 6, 12] as hour}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint('hour'), hour * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Previous</Badge>
							{hour} hour{hour > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="day">
					{#each [1, 7, 15] as day}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint('day'), day * 24 * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Previous</Badge>
							{day} day{day > 1 ? 's' : ''}
						</Command.Item>
					{/each}
				</Command.Group>
				<Command.Separator />
				<Command.Group heading="month">
					{#each [30, 45, 90, 180] as days}
						<Command.Item
							onclick={() =>
								setValue(startSetter, endSetter, getStartPoint('day'), days * 24 * 60 * 60 * 1000)}
						>
							<Badge variant="outline">Previous</Badge>
							{days} days
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Tabs.Content>
</Tabs.Root>
