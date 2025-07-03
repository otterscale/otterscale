<script lang="ts" module>
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		selectedVolume = $bindable()
	}: { selectedScope: string; selectedFacility: string; selectedVolume: string } = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const volumeOptions = writable<SingleSelect.OptionType[]>([]);
	let isVolumesLoading = $state(true);
	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listVolumes({
				scopeUuid: selectedScope,
				facilityName: selectedFacility
			});
			volumeOptions.set(
				response.volumes.map(
					(volume) =>
						({
							value: volume.name,
							label: volume.name,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isVolumesLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchVolumeOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});
</script>

{#if isMounted}
	<SingleSelect.Root options={volumeOptions} bind:value={selectedVolume}>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>No results found.</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $volumeOptions as option}
							<SingleSelect.Item
								{option}
								onclick={() => {
									selectedScope = option.value.scopeUuid;
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
