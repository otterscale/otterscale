<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
	import {
		Nexus,
		type CreateMachineRequest,
		type Scope,
		type Machine,
		type Machine_Placement,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const tagsStore = writable<Tag[]>();
	const tagsLoading = writable(true);
	async function fetchTags() {
		try {
			const response = await client.listTags({});
			tagsStore.set(response.tags);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			tagsLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchTags();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});

	const DEFAULT_REQUEST = {
		id: machine.id,
		enableSsh: true,
		skipBmcConfig: true,
		skipNetworking: true,
		skipStorage: true,
		tags: [] as string[]
	} as CreateMachineRequest;

	let createMachineRequest = $state(DEFAULT_REQUEST);

	function reset() {
		createMachineRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:compass" />
		Add Machine
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add {machine.fqdn}</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="flex flex-col gap-3">
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<legent>Scope</legent>
						<Select.Root type="single" bind:value={createMachineRequest.scopeUuid}>
							<Select.Trigger>
								{createMachineRequest.scopeUuid ? createMachineRequest.scopeUuid : 'Select'}
							</Select.Trigger>
							<Select.Content class="w-fit">
								{#each $scopesStore as scope}
									<Select.Item value={scope.uuid}>
										{scope.name}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</fieldset>

					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<span class="flex justify-between">
							<legent>Enable SSH</legent>
							<Switch bind:checked={createMachineRequest.enableSsh} />
						</span>
					</fieldset>
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<span class="flex justify-between">
							<legent>Skip BMC Configuration</legent>
							<Switch id="skip_bmc_config" bind:checked={createMachineRequest.skipBmcConfig} />
						</span>
					</fieldset>
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<span class="flex justify-between">
							<legent>Skip Networking</legent>
							<Switch id="skip_networking" bind:checked={createMachineRequest.skipNetworking} />
						</span>
					</fieldset>
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<span class="flex justify-between">
							<legent>Skip Storage</legent>
							<Switch id="skip_storage" bind:checked={createMachineRequest.skipStorage} />
						</span>
					</fieldset>

					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<legent>Tags</legent>
						<div class="flex flex-wrap gap-1">
							{#each createMachineRequest.tags as tag}
								<Badge variant="secondary">
									{tag}
								</Badge>
							{/each}
						</div>
						<Select.Root type="multiple" bind:value={createMachineRequest.tags}>
							<Select.Trigger>Select</Select.Trigger>
							<Select.Content class="w-fit">
								{#each $tagsStore as tag}
									<Select.Item value={tag.name}>
										{tag.name}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</fieldset>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.createMachine(createMachineRequest).then((r) => {
						toast.info(`Create ${machine.fqdn}`);
					});
					// toast.info(`Create ${machine.fqdn}`);
					console.log(createMachineRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
