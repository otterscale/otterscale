<script lang="ts" module>
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	const modes = writable<SingleSelect.OptionType[]>([
		{ value: 'get', label: 'Get', icon: 'ph:download-simple' },
		{ value: 'put', label: 'Put', icon: 'ph:upload-simple' },
		{ value: 'delete', label: 'Delete', icon: 'ph:trash' }
	]);
	let { selectedMode = $bindable() }: { selectedMode: string } = $props();
</script>

<SingleSelect.Root options={modes} bind:value={selectedMode}>
	<SingleSelect.Trigger />
	<SingleSelect.Content>
		<SingleSelect.Options>
			<SingleSelect.Input />
			<SingleSelect.List>
				<SingleSelect.Empty>No results found.</SingleSelect.Empty>
				<SingleSelect.Group>
					{#each $modes as option}
						<SingleSelect.Item {option}>
							<Icon
								icon={option.icon ? option.icon : 'ph:empty'}
								class={cn('size-5', option.icon ? 'visible' : 'invisible')}
							/>
							{option.label}
							<SingleSelect.Check {option} />
						</SingleSelect.Item>
					{/each}
				</SingleSelect.Group>
			</SingleSelect.List>
		</SingleSelect.Options>
	</SingleSelect.Content>
</SingleSelect.Root>
