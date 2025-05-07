<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		Nexus,
		type CreateFacilityRequest,
		type Machine,
		type Machine_Constraint,
		type Machine_Placement,
		type Scope,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';

	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Tabs from '$lib/components/ui/tabs';

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

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	async function fetchFacilities(scopeUuid: string) {
		try {
			const response = await client.listFacilities({
				scopeUuid: scopeUuid
			});

			return response.facilities;
		} catch (error) {
			console.error('Error fetching:', error);
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

	const DEFAULT_PLACEMENTS = [] as Machine_Placement[];
	const DEFAULT_CONSTRAINT = { tags: [] as string[] } as Machine_Constraint;
	const DEFAULT_REQUEST = {
		placements: DEFAULT_PLACEMENTS,
		constraint: DEFAULT_CONSTRAINT,
		trust: true
	} as CreateFacilityRequest;

	let createFacilityRequest = $state(DEFAULT_REQUEST);
	let createFacilityPlacements = $state(DEFAULT_PLACEMENTS);
	let createFacilityConstraint = $state(DEFAULT_CONSTRAINT);

	function reset() {
		createFacilityRequest = DEFAULT_REQUEST;
		resetCreateFacilityPlacement();
		resetCreateFacilityConstraint();
	}
	function resetCreateFacilityPlacement() {
		createFacilityPlacements = DEFAULT_PLACEMENTS;
	}
	function resetCreateFacilityConstraint() {
		createFacilityConstraint = DEFAULT_CONSTRAINT;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	let resourceMachineActive = $derived(createFacilityPlacements.length > 0);
	let resourceConstraintActive = $derived(
		Boolean(
			createFacilityConstraint.architecture ||
				createFacilityConstraint.cpuCores ||
				createFacilityConstraint.memoryMb ||
				createFacilityConstraint.tags?.length
		)
	);

	let scopeUuidToName = $derived(
		$scopesStore.reduce(
			(m, scope) => {
				m[scope.uuid] = scope.name;
				return m;
			},
			{} as Record<string, string>
		)
	);

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchMachines();
			await fetchTags();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>Install</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Add Facility</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid max-h-[77vh] w-full gap-3 overflow-y-auto">
					<Tabs.Root value="basic_information">
						<Tabs.List class="w-fit">
							<Tabs.Trigger value="basic_information">Basic Information</Tabs.Trigger>
							<Tabs.Trigger value="resource">Resource</Tabs.Trigger>
						</Tabs.List>
						<Tabs.Content value="basic_information">
							<div class="grid items-center gap-3 p-3">
								<Label class="text-sm font-medium">Scope</Label>
								<Select.Root type="single" bind:value={createFacilityRequest.scopeUuid}>
									<Select.Trigger>
										{scopeUuidToName[createFacilityRequest.scopeUuid] ?? 'Select'}
									</Select.Trigger>
									<Select.Content>
										{#each $scopesStore as scope}
											<Select.Item value={scope.uuid}>{scope.name}</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>

								<Label class="text-sm font-medium">Facility</Label>
								<!-- <Input bind:value={createFacilityRequest.name} /> -->
								{#await fetchFacilities(createFacilityRequest.scopeUuid)}
									<Input placeholder={'Loading...'} disabled />
								{:then facilitiesByScope}
									{#if facilitiesByScope && facilitiesByScope.length > 0}
										<Select.Root type="single" bind:value={createFacilityRequest.name}>
											<Select.Trigger>
												{createFacilityRequest.name ? createFacilityRequest.name : 'Select'}
											</Select.Trigger>
											<Select.Content>
												{#each facilitiesByScope || [] as facility}
													<Select.Item value={facility.name}>{facility.name}</Select.Item>
												{/each}
											</Select.Content>
										</Select.Root>
									{:else}
										<Input placeholder={'No facility'} disabled />
									{/if}
								{/await}

								<Label class="text-sm font-medium">Channel</Label>
								<Input bind:value={createFacilityRequest.channel} />
								<Label class="text-sm font-medium">Revision</Label>
								<Input bind:value={createFacilityRequest.revision} type="number" />
								<Label class="text-sm font-medium">Number</Label>
								<Input bind:value={createFacilityRequest.number} type="number" />
								<Label class="text-sm font-medium">Charm</Label>
								<Input bind:value={createFacilityRequest.charmName} />
								<Label class="text-sm font-medium">Configuration</Label>
								<Input bind:value={createFacilityRequest.configYaml} />
								<div class="flex justify-between">
									<Label class="text-sm font-medium">Trust</Label>
									<Switch bind:checked={createFacilityRequest.trust} />
								</div>
							</div>
						</Tabs.Content>
						<Tabs.Content value="resource">
							<fieldset
								class="grid items-center gap-3 rounded-lg border p-3"
								class:opacity-50={resourceConstraintActive}
							>
								<legend class="text-sm font-semibold">Machine</legend>
								<DropdownMenu.Root>
									<DropdownMenu.Trigger class="w-full">
										<Button
											variant="outline"
											class="w-full justify-between"
											disabled={resourceConstraintActive}
										>
											<span>Select Machines</span>
											<Icon icon="lucide:chevron-down" />
										</Button>
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="end" class="max-h-[300px] overflow-auto">
										{#each $machinesStore as machine}
											<DropdownMenu.Sub>
												<DropdownMenu.SubTrigger>{machine.fqdn}</DropdownMenu.SubTrigger>
												<DropdownMenu.SubContent>
													{#each ['lxd', 'kvm', 'machine'] as type}
														<DropdownMenu.Item
															onclick={() => {
																createFacilityPlacements = [
																	...createFacilityPlacements,
																	{
																		machineId: machine.id,
																		type: {
																			case: type,
																			value: true
																		}
																	} as Machine_Placement
																];
															}}
														>
															{type}
														</DropdownMenu.Item>
													{/each}
												</DropdownMenu.SubContent>
											</DropdownMenu.Sub>
										{/each}
									</DropdownMenu.Content>
								</DropdownMenu.Root>
								<div class="flex flex-wrap gap-2">
									{#each createFacilityPlacements as placement}
										<Badge
											variant="outline"
											onclick={() => {
												createFacilityPlacements = createFacilityPlacements.filter(
													(p) => p.machineId !== placement.machineId
												);
											}}
										>
											<div class="flex w-fit items-center gap-1">
												{placement.machineId}: {placement.type.case}
												<Icon icon="ph:x" />
											</div>
										</Badge>
									{/each}
								</div>
							</fieldset>

							<div class="relative">
								<div class="absolute inset-0 flex items-center">
									<span class="w-full border-t"></span>
								</div>
								<div class="relative flex justify-center text-xs">
									<span class="bg-background p-3 text-muted-foreground">Or</span>
								</div>
							</div>

							<fieldset
								class="grid items-center gap-3 rounded-lg border p-3"
								class:opacity-50={resourceMachineActive}
							>
								<legend class="text-sm font-semibold">Constraint</legend>

								<Label class="text-sm font-medium">Architecture</Label>
								<Input
									bind:value={createFacilityConstraint.architecture}
									disabled={resourceMachineActive}
								/>
								<Label class="text-sm font-medium">CPU Cores</Label>
								<Input
									type="number"
									bind:value={createFacilityConstraint.cpuCores}
									disabled={resourceMachineActive}
								/>
								<Label class="text-sm font-medium">Memory (MB)</Label>
								<Input
									type="number"
									bind:value={createFacilityConstraint.memoryMb}
									disabled={resourceMachineActive}
								/>
								<Label class="text-sm font-medium">Tags</Label>
								<span class="flex flex-wrap gap-3">
									{#each createFacilityConstraint.tags || [] as tag}
										<Badge variant="outline">
											<div class="flex w-fit items-center gap-1">
												{tag}
												<Button
													variant="ghost"
													class="h-4 w-4 p-0 hover:text-destructive"
													disabled={resourceMachineActive}
													onclick={() => {
														createFacilityConstraint.tags =
															createFacilityConstraint.tags?.filter((t) => t !== tag) || [];
													}}
												>
													<Icon icon="ph:x" />
												</Button>
											</div>
										</Badge>
									{/each}
									<Select.Root
										type="multiple"
										bind:value={createFacilityConstraint.tags}
										disabled={resourceMachineActive}
									>
										<Select.Trigger>Select</Select.Trigger>
										<Select.Content>
											{#each $tagsStore || [] as tag}
												<Select.Item value={tag.name}>
													{tag.name}
												</Select.Item>
											{/each}
										</Select.Content>
									</Select.Root>
								</span>
							</fieldset>
						</Tabs.Content>
					</Tabs.Root>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.createFacility(createFacilityRequest).then((r) => {
						toast.success(`Create ${createFacilityRequest.name}.`);
					});
					console.log(createFacilityRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
