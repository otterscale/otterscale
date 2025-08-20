<script lang="ts">
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';

	const modes = writable<SingleSelect.OptionType[]>([
		{ value: 'read', label: 'Read', icon: 'ph:download-simple' },
		{ value: 'write', label: 'Write', icon: 'ph:upload-simple' },
		{ value: 'trim', label: 'Trim', icon: 'ph:broom' }
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
