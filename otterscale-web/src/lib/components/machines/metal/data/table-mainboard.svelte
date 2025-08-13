<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Table from '$lib/components/custom/table';
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
			<Table.Head>VENDOR</Table.Head>
			<Table.Head>PRODUCT</Table.Head>
			<Table.Head>FIRMWARE</Table.Head>
			<Table.Head>BOOT MODE</Table.Head>
			<Table.Head>VERSION</Table.Head>
			<Table.Head>DATE</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		<Table.Row class="*:text-xs">
			<Table.Cell>{$machine.hardwareInformation.mainboard_vendor}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_product}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_vendor}</Table.Cell>
			<Table.Cell>{$machine.biosBootMethod.toUpperCase()}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_version}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.mainboard_firmware_date}</Table.Cell>
		</Table.Row>
	</Table.Body>
</Table.Root>
