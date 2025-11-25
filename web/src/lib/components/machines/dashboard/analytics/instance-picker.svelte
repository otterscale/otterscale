<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';

	import { MachineService } from '$lib/api/machine/v1/machine_pb';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { selectedInstance = $bindable() }: { selectedInstance: string | undefined } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const instanceOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	async function fetchOptions() {
		try {
			const response = await machineClient.listMachines({});
			instanceOptions.set([
				...response.machines
					.filter((machine) =>
						machine.workloadAnnotations?.['juju-machine-id']?.includes('-machine-')
					)
					.map((machine) => ({
						value: machine.fqdn,
						label: machine.fqdn,
						icon: 'ph:desktop'
					})),
				{
					value: '.*',
					label: 'All Machines',
					icon: 'ph:desktop-duotone'
				}
			]);
			selectedInstance = $instanceOptions[0].value;
		} catch (error) {
			console.error('Failed to fetch machines:', error);
		}
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetchOptions();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<SingleSelect.Root options={instanceOptions} bind:value={selectedInstance}>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $instanceOptions as option (option.value)}
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
{/if}
