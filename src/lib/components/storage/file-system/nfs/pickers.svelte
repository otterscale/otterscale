<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		scope,
		volume,
		selectedSubvolumeGroupName = $bindable()
	}: {
		scope: string;
		volume: string;
		selectedSubvolumeGroupName: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const subvolumeGroupOptions = writable<SingleSelect.OptionType[]>([]);
	async function fetchVolumeOptions() {
		try {
			const response = await storageClient.listSubvolumeGroups({
				scope: scope,
				volumeName: volume
			});

			subvolumeGroupOptions.set(
				response.subvolumeGroups.map(
					(subvolumeGroup) =>
						({
							value: subvolumeGroup.name,
							label: subvolumeGroup.name,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);
			subvolumeGroupOptions.update((origin) => [
				...origin,
				{
					value: '',
					label: 'default',
					icon: 'ph:cube'
				} as SingleSelect.OptionType
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

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>{m.group()}</Picker.Label>
		{#if isMounted}
			<SingleSelect.Root options={subvolumeGroupOptions} bind:value={selectedSubvolumeGroupName}>
				<SingleSelect.Trigger />
				<SingleSelect.Content>
					<SingleSelect.Options>
						<SingleSelect.Input />
						<SingleSelect.List>
							<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
							<SingleSelect.Group>
								{#each $subvolumeGroupOptions as option (option.value)}
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
		{:else}
			<Loading.Selection />
		{/if}
	</Picker.Wrapper>
</Picker.Root>
