<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import { m } from '$lib/paraglide/messages';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		machine,
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>
				{m.name()}
				<Table.SubHead>{m.mac_address()}</Table.SubHead>
			</Table.Head>
			<Table.Head>
				{m.ip_address()}
				<Table.SubHead>{m.subnet()}</Table.SubHead>
			</Table.Head>
			<Table.Head>
				{m.link_speed()}
				<Table.SubHead>{m.link_connected()}</Table.SubHead>
			</Table.Head>
			<Table.Head>
				{m.fabric()}
				<Table.SubHead>{m.vlan()}</Table.SubHead>
			</Table.Head>
			<Table.Head>{m.type()}</Table.Head>
			<Table.Head>{m.dhcp_on()}</Table.Head>
			<Table.Head>{m.boot_interface()}</Table.Head>
			<Table.Head>{m.interface_speed()}</Table.Head>
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
