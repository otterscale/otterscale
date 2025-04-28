<script lang="ts">
	// External dependencies
	import Icon from '@iconify/svelte';
	import { capitalizeFirstLetter } from 'better-auth';
	import { toast } from 'svelte-sonner';

	// Internal UI components
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Switch } from '$lib/components/ui/switch/index.js';

	// Internal utilities and types
	import { type Machine } from '$gen/api/stack/v1/stack_pb';
	import { nodeIcon } from '$lib/node';
	import ListExpression from './ui/list-expression.svelte';
	import ChartMetric from './ui/chart-metric.svelte';
	import CardDictionary from './ui/card-dictionary.svelte';
	import {
		type ImportBootResourcesRequest,
		type PowerOnMachineRequest,
		type PowerOffMachineRequest,
		type CommissionMachineRequest
	} from '$gen/api/stack/v1/stack_pb';

	const nodeType = 'MAAS';

	let commissionMachine = $state({} as CommissionMachineRequest);

	let {
		machine
	}: {
		machine: Machine;
	} = $props();
</script>

{#snippet NetworkHeader(header: string, subheader?: string)}
	<div class="text-sm">
		{header}
		{#if subheader}
			<div class="text-xs font-light text-muted-foreground">
				{subheader}
			</div>
		{/if}
	</div>
{/snippet}

<div class="grid gap-3 px-8">
	<div class="flex items-center">
		<div class="flex items-center space-x-2">
			<Icon icon={nodeIcon(nodeType)} class="size-8" />
			<div class="flex-col p-2">
				<div class="font-bold">{machine.fqdn}</div>
				<div class="flex text-sm text-muted-foreground">
					{machine.systemId}
				</div>
			</div>
		</div>
		<div class="ml-auto">
			<Badge
				variant="secondary"
				class="cursor-pointer hover:scale-105"
				onclick={() => {
					if (confirm('Delete this item?')) {
						toast.success(`Item ID ${machine.systemId} deleted.`);
					}
				}}
			>
				<Icon icon="ph:trash" class="size-5 sm:flex" />
			</Badge>
		</div>
	</div>
	<div class="space-y-1">
		<p>{machine.description}</p>
		<div>
			<AlertDialog.Root>
				<AlertDialog.Trigger>
					<Badge variant="outline" class="cursor-pointer text-muted-foreground hover:scale-105">
						<Icon
							icon="ph:power"
							class="size-5 sm:flex"
							style="color: {machine.powerState === 'on' ? 'green' : 'inherit'}"
						/>
						<span class="pl-1 text-sm">Power {machine.powerState}</span>
					</Badge>
				</AlertDialog.Trigger>
				<AlertDialog.Content>
					<AlertDialog.Header>
						<AlertDialog.Title
							>Power {machine.powerState === 'on' ? 'off' : 'on'}
							{machine.fqdn}</AlertDialog.Title
						>
						<AlertDialog.Description>
							Are you sure you want to {machine.powerState === 'on' ? 'turn off' : 'turn on'} this machine?
						</AlertDialog.Description>
					</AlertDialog.Header>
					<AlertDialog.Footer>
						<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
						<AlertDialog.Action
							onclick={() => {
								toast.promise(
									fetch(`/api/management/machine/${machine.systemId}/power`, {
										method: 'POST',
										body: JSON.stringify({
											action: machine.powerState === 'on' ? 'off' : 'on'
										})
									}),
									{
										loading: 'In Processing...',
										success: 'Done.',
										error: 'Fail!'
									}
								);
							}}
						>
							Confirm
						</AlertDialog.Action>
					</AlertDialog.Footer>
				</AlertDialog.Content>
			</AlertDialog.Root>
			<AlertDialog.Root>
				<AlertDialog.Trigger>
					<Badge variant="outline" class="cursor-pointer text-muted-foreground hover:scale-105">
						<Icon icon="ph:computer-tower" class="size-5 sm:flex" />
						<span class="pl-1 text-sm"></span>Commission
					</Badge>
				</AlertDialog.Trigger>
				<AlertDialog.Content>
					<AlertDialog.Header>
						<AlertDialog.Title>Commision {machine.fqdn}</AlertDialog.Title>
						<AlertDialog.Description>
							<fieldset class="grid gap-4 py-4">
								<div class="flex items-center space-x-2">
									<Switch id="enable_ssh" bind:checked={commissionMachine.enableSsh} />
									<label for="enable_ssh">Enable SSH</label>
								</div>
								<div class="flex items-center space-x-2">
									<Switch id="skip_bmc_config" bind:checked={commissionMachine.skipBmcConfig} />
									<label for="skip_bmc_config">Skip BMC Configuration</label>
								</div>
								<div class="flex items-center space-x-2">
									<Switch id="skip_networking" bind:checked={commissionMachine.skipNetworking} />
									<label for="skip_networking">Skip Networking</label>
								</div>
								<div class="flex items-center space-x-2">
									<Switch id="skip_storage" bind:checked={commissionMachine.skipStorage} />
									<label for="skip_storage">Skip Storage</label>
								</div>
							</fieldset>
						</AlertDialog.Description>
					</AlertDialog.Header>
					<AlertDialog.Footer>
						<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
						<AlertDialog.Action
							onclick={() => {
								toast.info(
									[
										commissionMachine.enableSsh,
										commissionMachine.skipBmcConfig,
										commissionMachine.skipNetworking,
										commissionMachine.skipStorage
									].join(', ')
								);
							}}
						>
							Confirm
						</AlertDialog.Action>
					</AlertDialog.Footer>
				</AlertDialog.Content>
			</AlertDialog.Root>
			<Badge
				variant="outline"
				class="cursor-pointer text-muted-foreground hover:scale-105"
				onclick={() => {
					toast.promise(
						fetch(`/api/management/machine/${machine.systemId}/boot-resources`, {
							method: 'POST'
						}),
						{
							loading: 'Importing boot resources...',
							success: 'Boot resources imported successfully',
							error: 'Failed to import boot resources'
						}
					);
				}}
			>
				<Icon icon="ph:download" class="size-5 sm:flex" />
				<span class="pl-1 text-sm">Import Boot Resources</span>
			</Badge>
		</div>
		<div class="space-x-1">
			{#each machine.tags as tag}
				<Badge variant="outline" class="text-muted-foreground">
					<Icon icon="ph:tag" class="size-5 sm:flex" />
					<span class="pl-1 text-sm">{tag}</span>
				</Badge>
			{/each}
		</div>
	</div>
	<Tabs.Root value="summary">
		<Tabs.List class={`grid w-full grid-cols-3`}>
			<Tabs.Trigger value="summary">Summary</Tabs.Trigger>
			<Tabs.Trigger value="hardware_information">Hardware Information</Tabs.Trigger>
			<Tabs.Trigger value="network">Networks</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="summary">
			<div class="grid gap-3">
				<div class="grid grid-cols-4 gap-3">
					<ChartMetric
						description={'VIRTUAL'}
						metric={machine.status}
						footer={`${capitalizeFirstLetter(machine.osystem)} ${machine.hweKernel} ${capitalizeFirstLetter(machine.distroSeries)}`}
					/>
					<ChartMetric
						description={'CPU'}
						metric={`${machine.cpuCount} cores`}
						footer={machine.hardwareInformation['cpu_model']}
						tag={machine.architecture}
					/>
					<ChartMetric
						description={'MEMORY'}
						metric={`${Math.round(Number(machine.memory) / 1024)} GiB`}
					/>
					<ChartMetric
						description={'STORAGE'}
						metric={`${Math.round(Number(machine.storage) / 1024)} GiB`}
					/>
				</div>
				<div class="grid grid-cols-1 gap-3">
					<CardDictionary title={'WORKLOAD ANNOTATIONS'}>
						<Table.Root>
							<Table.Body>
								<Table.Row>
									<Table.Cell>JUJU Controller UUID</Table.Cell>
									<Table.Cell>
										{machine.workloadAnnotations['juju-controller-uuid']}
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>JUJU Machine ID</Table.Cell>
									<Table.Cell>
										{machine.workloadAnnotations['juju-machine-id']}
									</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>JUJU Model UUID</Table.Cell>
									<Table.Cell>
										{machine.workloadAnnotations['juju-model-uuid']}
									</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					</CardDictionary>

					<!-- <Card.Root class="col-span-2">
				<Card.Header class="text-text-base">
					<Card.Title>NUMA NODES</Card.Title>
				</Card.Header>
				<Card.Content class="grid flex-grow gap-3 text-xl">
					{#each machine.numaNodes as numaNode}
						<p class="text-sm">Node {numaNode['index']}</p>
						<Table.Root>
							<Table.Body class="text-xs">
								<Table.Row>
									<Table.Cell>CPU cores</Table.Cell>
									<Table.Cell>{numaNode['cores']}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Memory</Table.Cell>
									<Table.Cell>{Math.round(Number(numaNode['memory']) / 1024)} GiB</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					{/each}
				</Card.Content>
			</Card.Root> -->
				</div>
			</div>
		</Tabs.Content>
		<Tabs.Content value="hardware_information">
			<div class="grid gap-3">
				<div class="grid grid-cols-3 gap-3">
					<CardDictionary title={'System'}>
						<Table.Root>
							<Table.Body>
								<Table.Row>
									<Table.Cell>Vendor</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_vendor}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Product</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_product}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Version</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_version}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Serial</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_serial}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>SKU</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_sku}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Family</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.system_family}</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					</CardDictionary>
					<CardDictionary title={'Mainboard'}>
						<Table.Root>
							<Table.Body>
								<Table.Row>
									<Table.Cell>Vendor</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.mainboard_vendor}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Product</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.mainboard_product}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Firmware</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.mainboard_firmware_vendor}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Boot Mode</Table.Cell>
									<Table.Cell>{machine.biosBootMethod.toUpperCase()}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Version</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.mainboard_firmware_version}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Date</Table.Cell>
									<Table.Cell>{machine.hardwareInformation.mainboard_firmware_date}</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					</CardDictionary>
					<CardDictionary title={'Chassis'}>
						<Table.Root>
							<Table.Body>
								<Table.Row>
									<Table.Cell>Vendor</Table.Cell>
									<Table.Cell>{machine.hardwareInformation['chassis_vendor']}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Type</Table.Cell>
									<Table.Cell>{machine.hardwareInformation['chassis_type']}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Version</Table.Cell>
									<Table.Cell>{machine.hardwareInformation['chassis_version']}</Table.Cell>
								</Table.Row>
								<Table.Row>
									<Table.Cell>Serial</Table.Cell>
									<Table.Cell>{machine.hardwareInformation['chassis_serial']}</Table.Cell>
								</Table.Row>
							</Table.Body>
						</Table.Root>
					</CardDictionary>
				</div>
				<CardDictionary title={'Block Devices'}>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Name</Table.Head>
								<Table.Head>Model</Table.Head>
								<Table.Head>Serial</Table.Head>
								<Table.Head>Boot Disk</Table.Head>
								<Table.Head>Firmware Version</Table.Head>
								<Table.Head>Size</Table.Head>
								<Table.Head>Type</Table.Head>
								<Table.Head>User For</Table.Head>
								<Table.Head>Tags</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each machine.blockDevices as blockDevice}
								<Table.Row>
									<Table.Cell>{blockDevice.name}</Table.Cell>
									<Table.Cell>{blockDevice.model}</Table.Cell>
									<Table.Cell>{blockDevice.serial}</Table.Cell>
									<Table.Cell>{blockDevice.bootDisk}</Table.Cell>
									<Table.Cell>{blockDevice.firmwareVersion}</Table.Cell>
									<Table.Cell>{blockDevice.size}</Table.Cell>
									<Table.Cell>{blockDevice.type}</Table.Cell>
									<Table.Cell>{blockDevice.usedFor}</Table.Cell>
									<Table.Cell><ListExpression items={blockDevice.tags} /></Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardDictionary>
			</div>
		</Tabs.Content>
		<Tabs.Content value="network">
			<CardDictionary>
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>{@render NetworkHeader('NAME', 'MAC Address')}</Table.Head>
							<Table.Head>{@render NetworkHeader('IP ADDRESS', 'Subnet')}</Table.Head>
							<Table.Head>{@render NetworkHeader('LINK SPEED', 'Link Connected')}</Table.Head>
							<Table.Head>{@render NetworkHeader('FABRIC', 'VLAN')}</Table.Head>
							<Table.Head>{@render NetworkHeader('Type')}</Table.Head>
							<Table.Head>{@render NetworkHeader('DHCP On')}</Table.Head>
							<Table.Head>{@render NetworkHeader('Boot Interface')}</Table.Head>
							<Table.Head>{@render NetworkHeader('Interface Speed')}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each machine.networkInterfaces as networkInterface}
							<Table.Row>
								<Table.Cell>
									{networkInterface.name}
									<div>
										{networkInterface.macAddress}
									</div>
								</Table.Cell>
								<Table.Cell>
									{networkInterface.ipAddress}
									<div>
										{networkInterface.subnetName}
									</div>
								</Table.Cell>
								<Table.Cell>
									{networkInterface.linkSpeed} Mbps
									<div>
										<Icon
											icon={networkInterface.linkConnected ? 'ph:check-circle' : 'ph:x-circle'}
											style="color: {networkInterface.linkConnected ? 'green' : 'red'}"
										/>
									</div>
								</Table.Cell>
								<Table.Cell>
									{networkInterface.fabricName}
									<div>
										{networkInterface.vlanName}
									</div>
								</Table.Cell>
								<Table.Cell>{networkInterface.type}</Table.Cell>
								<Table.Cell>
									<Icon
										icon={networkInterface.dhcpOn ? 'ph:check' : 'ph:x'}
										style="color: {networkInterface.dhcpOn ? 'green' : 'red'}"
									/>
								</Table.Cell>
								<Table.Cell>{networkInterface.bootInterface}</Table.Cell>
								<Table.Cell>{networkInterface.interfaceSpeed} Mbps</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</CardDictionary>
		</Tabs.Content>
	</Tabs.Root>
</div>
