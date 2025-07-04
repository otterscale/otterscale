<script lang="ts" module>
	import { EssentialService, Essential_Type } from '$gen/api/essential/v1/essential_pb';
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
		selectedScope = $bindable(),
		selectedFacility = $bindable()
	}: { selectedScope: string; selectedFacility: string } = $props();

	let selectedCeph = $state({});

	const transport: Transport = getContext('transport');
	const essentialClient = createClient(EssentialService, transport);

	const cephOptions = writable<SingleSelect.OptionType[]>([]);
	let isCephsLoading = $state(true);
	async function fetchCephOptions() {
		try {
			const response = await essentialClient.listEssentials({ type: Essential_Type.CEPH });
			cephOptions.set(
				response.essentials.map(
					(essential) =>
						({
							value: { scopeUuid: essential.scopeUuid, facilityName: essential.name },
							label: `${essential.scopeName}-${essential.name}`,
							icon: 'ph:cube'
						}) as SingleSelect.OptionType
				)
			);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isCephsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchCephOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		isMounted = true;
	});
</script>

{#if isMounted}
	<SingleSelect.Root options={cephOptions} bind:value={selectedCeph}>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>No results found.</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $cephOptions as option}
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
{:else}
	<Skeleton class="bg-muted w-[100px]" />
{/if}
