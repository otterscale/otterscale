<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import {
		ConfigurationService,
		type InternalObjectService,
		InternalObjectService_Type
	} from '$lib/api/configuration/v1/configuration_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		scope,
		selectedInternalObjectService = $bindable()
	}: { scope: string; selectedInternalObjectService: InternalObjectService } = $props();

	let selectedInit = $state({});

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);

	const internalObjectServices = writable<SingleSelect.OptionType[]>([]);
	async function fetchOptions() {
		try {
			const response = await client.listInternalObjectServices({
				scope: scope
			});
			internalObjectServices.set(
				response.internalObjectServices.map(
					(internalObjectService) =>
						({
							value: internalObjectService,
							label: `${InternalObjectService_Type[internalObjectService.type]}`,
							icon: 'ph:cube',
							information: `${InternalObjectService_Type[internalObjectService.type]} (${internalObjectService.host})`
						}) as SingleSelect.OptionType
				)
			);
			if (selectedInternalObjectService) {
				const options = response.internalObjectServices.map(
					(internalObjectService) =>
						({
							value: internalObjectService,
							label: `${InternalObjectService_Type[internalObjectService.type]}`,
							icon: 'ph:cube',
							information: `${InternalObjectService_Type[internalObjectService.type]} (${internalObjectService.host})`
						}) as SingleSelect.OptionType
				);
				const matched = options.find(
					(opt) =>
						opt.value.type === selectedInternalObjectService.type &&
						opt.value.scope === selectedInternalObjectService.scope &&
						opt.value.host === selectedInternalObjectService.host
				);
				if (matched) {
					selectedInit = matched.value;
				}
			}
		} catch (error) {
			console.error('Error fetching:', error);
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
	<SingleSelect.Root options={internalObjectServices} bind:value={selectedInit}>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>No results found.</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $internalObjectServices as option}
							<SingleSelect.Item
								{option}
								onclick={() => {
									selectedInternalObjectService = option.value;
								}}
							>
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
	<Skeleton class="w-[100px] bg-muted" />
{/if}
