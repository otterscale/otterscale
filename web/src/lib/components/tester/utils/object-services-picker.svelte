<script lang="ts" module>
	import { BISTService } from '$gen/api/bist/v1/bist_pb'
	import * as Picker from '$lib/components/custom/picker';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { cn } from '$lib/utils.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedName = $bindable(),
		selectedPath = $bindable()
	}: { selectedName: string; selectedPath: string } = $props();

	let tmp = $state({});

	const transport: Transport = getContext('transport');
	const bistClient = createClient(BISTService, transport);

	const objectServiceName = writable<SingleSelect.OptionType[]>([]);
	const objectServicePath = writable<SingleSelect.OptionType[]>([]);
	let isLoading = $state(true);
	async function fetchOptions() {
		try {
			const response = await bistClient.listObjectServices({ });
			objectServiceName.set(
				response.minios.map(
					(minios) =>
						({
							value: { name: minios.name },
							label: minios.name,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);
			objectServicePath.set(
				response.minios.map(
					(minios) =>
						({
							value: { endpoint: minios.endpoint },
							label: minios.endpoint,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});
</script>

{#if isMounted}
	<Picker.Root align="left">
		<!-- Name -->
		<Picker.Wrapper class="*:h-8">
			<SingleSelect.Root options={objectServiceName} bind:value={tmp}>
				<SingleSelect.Trigger />
				<SingleSelect.Content>
					<SingleSelect.Options>
						<SingleSelect.Input />
						<SingleSelect.List>
							<SingleSelect.Empty>No results found.</SingleSelect.Empty>
							<SingleSelect.Group>
								{#each $objectServiceName as option}
									<SingleSelect.Item
										{option}
										onclick={() => {
											selectedName = option.value;
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
			<!-- Path -->
			<SingleSelect.Root options={objectServicePath} bind:value={tmp}>
				<SingleSelect.Trigger />
				<SingleSelect.Content>
					<SingleSelect.Options>
						<SingleSelect.Input />
						<SingleSelect.List>
							<SingleSelect.Empty>No results found.</SingleSelect.Empty>
							<SingleSelect.Group>
								{#each $objectServicePath as option}
									<SingleSelect.Item
										{option}
										onclick={() => {
											selectedName = option.value;
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
		</Picker.Wrapper>
	</Picker.Root>
{:else}
	<Skeleton class="bg-muted w-[100px]" />
{/if}
