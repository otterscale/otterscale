<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Card from '$lib/components/ui/card';
	import Label from '$lib/components/ui/label/label.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<!-- <Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>{m.vendor()}</Table.Head>
			<Table.Head>{m.product()}</Table.Head>
			<Table.Head>{m.version()}</Table.Head>
			<Table.Head>{m.serial()}</Table.Head>
			<Table.Head>{m.sku()}</Table.Head>
			<Table.Head>{m.family()}</Table.Head>
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
</Table.Root> -->

<Card.Root class="w-full max-w-sm shadow-none">
	<Card.Header class="-mt-6 rounded-t-lg bg-muted p-2">
		<Card.Title class="text-center">{m.system()}</Card.Title>
	</Card.Header>
	<Card.Content>
		<div class="flex flex-col gap-6 text-sm">
			<div class="flex gap-2">
				<Icon icon="ph:cpu" />
				<div>
					<Label>{m.vendor()}</Label>
					<p>{$machine.hardwareInformation.system_vendor}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:archive-box" />
				<div>
					<Label>{m.product()}</Label>
					<p>{$machine.hardwareInformation.system_product}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:tag" />
				<div>
					<Label>{m.version()}</Label>
					<p>{$machine.hardwareInformation.system_version}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:barcode" />
				<div>
					<Label>{m.serial()}</Label>
					<p>{$machine.hardwareInformation.system_serial}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:ticket" />
				<div>
					<Label>{m.sku()}</Label>
					<p>{$machine.hardwareInformation.system_sku}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Icon icon="ph:users-three" />
				<div>
					<Label>{m.family()}</Label>
					<p>{$machine.hardwareInformation.system_family}</p>
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>
