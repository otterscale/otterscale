<script lang="ts">
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { writable, type Writable } from 'svelte/store';
	import { fetchSubvolumeGroupList, fetchSubvolumeGroup } from './data';
	import { DataTable } from './data-table';
	import Icon from '@iconify/svelte';
	import { cn } from '$lib/utils.js';
	import Label from '$lib/components/ui/label/label.svelte';

	function camelToPascal(value: string): string {
		return value
			.replace(/([a-z])([A-Z])/g, '$1 $2')
			.split(' ')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	let group = $state('');

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
	const data = $derived(fetchSubvolumeGroup(group));
</script>

<div class="flex items-center justify-end gap-4">
	<span class="flex items-center gap-2">
		<Label class="bg-muted h-8 whitespace-nowrap rounded-lg p-2 text-center">Group</Label>
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
	</span>
</div>
{#if data}
	<DataTable {data} />
{:else}
	<PageLoading />
{/if}
