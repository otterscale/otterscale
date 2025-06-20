<script lang="ts">
	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import { fetchSubvolumeGroup, fetchSubvolumeGroupList } from './utils.svelte';

	function camelToPascal(value: string): string {
		return value
			.replace(/([a-z])([A-Z])/g, '$1 $2')
			.split(' ')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	let { group = $bindable() }: { group: string } = $props();

	const groups = fetchSubvolumeGroupList();
	const groupOptions: Writable<SingleSelect.OptionType[]> = writable(
		groups.map(
			(group) =>
				({
					value: group,
					label: group ? camelToPascal(group) : 'Default',
					icon: 'ph:cube'
				}) as SingleSelect.OptionType
		)
	);
	[group] = groups;
</script>

<Picker.Root align="right">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Group</Picker.Label>
		<SingleSelect.Root options={groupOptions} bind:value={group}>
			<SingleSelect.Trigger />
			<SingleSelect.Content>
				<SingleSelect.Options>
					<SingleSelect.Input />
					<SingleSelect.List>
						<SingleSelect.Empty>No results found.</SingleSelect.Empty>
						<SingleSelect.Group>
							{#each $groupOptions as option}
								<SingleSelect.Item {option}>
									<Icon
										icon={option.icon ? option.icon : 'ph:empty'}
										class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
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
	</Picker.Wrapper>
</Picker.Root>
