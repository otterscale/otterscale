<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>
				NAME
				<Table.SubHead>MAC Address</Table.SubHead>
			</Table.Head>
			<Table.Head>
				IP ADDRESS
				<Table.SubHead>Subnet</Table.SubHead>
			</Table.Head>
			<Table.Head>
				LINK SPEED
				<Table.SubHead>Link Connected</Table.SubHead>
			</Table.Head>
			<Table.Head>
				FABRIC
				<Table.SubHead>VLAN</Table.SubHead>
			</Table.Head>
			<Table.Head>TYPE</Table.Head>
			<Table.Head>DHCP ON</Table.Head>
			<Table.Head>BOOT INTERFACE</Table.Head>
			<Table.Head>INTERFACE SPEED</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $machine.networkInterfaces as networkInterface}
			<Table.Row>
				<Table.Cell>
					{networkInterface.name}
					<Table.SubCell>
						{networkInterface.macAddress}
					</Table.SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.ipAddress}
					<Table.SubCell>
						{networkInterface.subnetName}
					</Table.SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.linkSpeed} Mbps
					<Table.SubCell>
						<Icon
							icon={networkInterface.linkConnected ? 'ph:check-circle' : 'ph:x-circle'}
							class={networkInterface.linkConnected ? 'text-primary' : 'text-destructive'}
						/>
					</Table.SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.fabricName}
					<Table.SubCell>
						{networkInterface.vlanName}
					</Table.SubCell>
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
