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
			<Table.Head>VERSION</Table.Head>
			<Table.Head>SERIAL</Table.Head>
			<Table.Head>SKU</Table.Head>
			<Table.Head>FAMILY</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		<Table.Row>
			<Table.Cell>{$machine.hardwareInformation.system_vendor}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.system_product}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.system_version}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.system_serial}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.system_sku}</Table.Cell>
			<Table.Cell>{$machine.hardwareInformation.system_family}</Table.Cell>
		</Table.Row>
	</Table.Body>
</Table.Root>
