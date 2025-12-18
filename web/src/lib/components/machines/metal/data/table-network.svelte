<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { SubCell, SubHead } from '$lib/components/custom/table';
	import * as Table from '$lib/components/ui/table';
	import { m } from '$lib/paraglide/messages';
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
				{m.name()}
				<SubHead>{m.mac_address()}</SubHead>
			</Table.Head>
			<Table.Head>
				{m.ip_address()}
				<SubHead>{m.subnet()}</SubHead>
			</Table.Head>
			<Table.Head>
				{m.link_speed()}
				<SubHead>{m.link_connected()}</SubHead>
			</Table.Head>
			<Table.Head>
				{m.fabric()}
				<SubHead>{m.vlan()}</SubHead>
			</Table.Head>
			<Table.Head>{m.type()}</Table.Head>
			<Table.Head>{m.dhcp_on()}</Table.Head>
			<Table.Head>{m.boot_interface()}</Table.Head>
			<Table.Head>{m.interface_speed()}</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each $machine.networkInterfaces as networkInterface, index (index)}
			<Table.Row>
				<Table.Cell>
					{networkInterface.name}
					<SubCell>
						{networkInterface.macAddress}
					</SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.ipAddress}
					<SubCell>
						{networkInterface.subnetName}
					</SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.linkSpeed} Mbps
					<SubCell>
						<Icon
							icon={networkInterface.linkConnected ? 'ph:check-circle' : 'ph:x-circle'}
							class={networkInterface.linkConnected ? 'text-primary' : 'text-destructive'}
						/>
					</SubCell>
				</Table.Cell>
				<Table.Cell>
					{networkInterface.fabricName}
					<SubCell>
						{networkInterface.vlanName}
					</SubCell>
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
