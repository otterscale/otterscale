<script lang="ts">
	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';
	import { fetchSubvolumeGroupList, fetchSubvolumeListByGroup } from './utils.svelte';

	function camelToPascal(value: string): string {
		return value
			.replace(/([a-z])([A-Z])/g, '$1 $2')
			.split(' ')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	let { group = $bindable(), volume = $bindable() }: { group: string; volume: string } = $props();

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

	const volumes = $derived(fetchSubvolumeListByGroup(group));
	const volumesOptions: Writable<SingleSelect.OptionType[]> = $derived(
		writable(
			volumes.map(
				(volume) =>
					({
						value: volume,
						label: camelToPascal(volume),
						icon: 'ph:cube'
					}) as SingleSelect.OptionType
			)
		)
	);

	[group] = groups;
	$effect(() => {
		[volume] = volumes;
	});
</script>

<Picker.Root align="right">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Group</Picker.Label>
		<SingleSelect.Root options={groupOptions} bind:value={group}>
			<SingleSelect.Trigger class="col-start-2 row-start-1 h-8 justify-start ring-0" />
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
	{#key group}
		<Picker.Wrapper class="*:h-8">
			<Picker.Label>Volume</Picker.Label>
			<SingleSelect.Root options={volumesOptions} bind:value={volume}>
				<SingleSelect.Trigger class="col-start-2 row-start-2 h-8 justify-start ring-0" />
				<SingleSelect.Content>
					<SingleSelect.Options>
						<SingleSelect.Input />
						<SingleSelect.List>
							<SingleSelect.Empty>No results found.</SingleSelect.Empty>
							<SingleSelect.Group>
								{#each $volumesOptions as option}
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
	{/key}
</Picker.Root>
