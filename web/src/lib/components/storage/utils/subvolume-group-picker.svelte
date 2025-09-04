<script lang="ts" module>
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroupName = $bindable(),
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const subvolumeGroupOptions = writable<SingleSelect.OptionType[]>([]);
	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listSubvolumeGroups({
				scopeUuid: selectedScopeUuid,
				facilityName: selectedFacility,
				volumeName: selectedVolume,
			});

			subvolumeGroupOptions.set(
				response.subvolumeGroups.map(
					(subvolumeGroup) =>
						({
							value: subvolumeGroup.name,
							label: subvolumeGroup.name,
							icon: 'ph:cube',
						}) as SingleSelect.OptionType,
				),
			);
			subvolumeGroupOptions.update((origin) => [
				...origin,
				{
					value: '',
					label: 'default',
					icon: 'ph:cube',
				} as SingleSelect.OptionType,
			]);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchVolumeOptions();
			isMounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isMounted}
	<SingleSelect.Root options={subvolumeGroupOptions} bind:value={selectedSubvolumeGroupName}>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>No results found.</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $subvolumeGroupOptions as option}
							<SingleSelect.Item
								{option}
								onclick={() => {
									selectedScopeUuid = option.value.scopeUuid;
									selectedFacility = option.value.facilityName;
								}}
							>
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
{/if}
