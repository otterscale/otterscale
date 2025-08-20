<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
	import { type Writable } from 'svelte/store';

	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>VENDOR</Table.Head>
			<Table.Head>TYPE</Table.Head>
			<Table.Head>VERSION</Table.Head>
			<Table.Head>SERIAL</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		<Table.Row>
			<Table.Cell>{$machine.hardwareInformation.chassis_vendor}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_type}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_version}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.chassis_serial}</Table.Cell>
		</Table.Row>
	</Table.Body>
</Table.Root>
