<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Drawer from '$lib/components/ui/drawer';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Accordion from '$lib/components/ui/accordion';
	import { formatTimeAgo } from '$lib/formatter';
	import pb from '$lib/pb';
	import * as Table from '$lib/components/ui/table';
	import * as Dialog from '$lib/components/ui/dialog';

	import { useId } from 'bits-ui';

	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Carousel from '$lib/components/ui/carousel';
	import * as Card from '$lib/components/ui/card';

	import * as Popover from '$lib/components/ui/popover/index.js';
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import { tick } from 'svelte';
	import * as Command from '$lib/components/ui/command';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Label } from '$lib/components/ui/label/index.js';
	import { FabricCreateOverview, FabricCreateSource } from '$lib/components/fabric';
	import FabricCreateDestination from '$lib/components/fabric/fabric-create-destination.svelte';
	import FabricCreateConnection from '$lib/components/fabric/fabric-create-connection.svelte';

	let items = $state([
		{
			name: 'Source',
			icon: 'ph:plug',
			active: false,
			set: {
				steps: [false, false],
				connectors: [
					{
						name: 'PostgreSQL',
						icon: 'logos:postgresql',
						parameters: [
							{
								key: 'connection_string',
								name: 'Connection String',
								value: '',
								description: `connection string, such as 'postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10'`
							},
							{
								key: 'namespace',
								name: 'Name Space',
								value: '',
								description: `namespace of database, such as 'public'`
							}
						],
						templates: [
							{
								name: 'TEST',
								parameters: [
									{
										key: 'connection_string',
										value: ''
									},
									{
										key: 'namespace',
										value: ''
									}
								]
							}
						]
					},
					{
						name: 'CSV',
						icon: 'ph:file-csv',
						parameters: [
							{
								key: 'file_path',
								name: 'File Path',
								value: '',
								description: 'csv file path'
							},
							{
								key: 'table_name',
								name: 'Table Name',
								value: '',
								description: 'destination table name'
							},
							{
								key: 'inferring',
								name: 'Inferring',
								value: '',
								description: ''
							}
						]
					}
				]
			}
		}
		// {
		// 	name: 'Destination',
		// 	icon: 'ph:plugs',
		// 	active: false
		// },
		// {
		// 	name: 'Connection',
		// 	icon: 'ph:plugs-connected',
		// 	active: false
		// },
		// {
		// 	name: 'Transformer',
		// 	icon: 'ph:tree-structure',
		// 	active: false
		// }
	]);

	let srcOpens: boolean[] = $state([false, false, false, false]);
	let dstOpens: boolean[] = $state([false, false, false, false]);
	let connectionOpens: boolean[] = $state([false, false, false, false]);
	let connectorOpens: boolean[] = $state([false, false, false, false]);

	let open = $state(false);
	//

	type Status = {
		value: string;
		label: string;
	};

	const statuses: Status[] = [
		{
			value: 'backlog',
			label: 'Backlog'
		},
		{
			value: 'todo',
			label: 'Todo'
		},
		{
			value: 'in progress',
			label: 'In Progress'
		},
		{
			value: 'done',
			label: 'Done'
		},
		{
			value: 'canceled',
			label: 'Canceled'
		}
	];

	let value = $state('');

	const selectedStatus = $derived(statuses.find((s) => s.value === value) ?? null);

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	const triggerId = useId();
</script>

<Drawer.Root direction="right">
	<Drawer.Trigger class={buttonVariants({ variant: 'outline' })}>Open Drawer</Drawer.Trigger>
	<Drawer.Content class="absolute inset-x-auto inset-y-0 right-0 w-3/5 space-y-2 rounded-tr-none">
		<Drawer.Header class="px-8 pt-0">
			<Drawer.Title class="flex items-center">
				<div class="flex items-center space-x-2">
					<Icon icon="logos:postgresql" class="size-8" />
					<div class="flex-col p-2">
						PostgreSQL
						<div class="flex items-center gap-1 text-sm text-muted-foreground">
							{pb.authStore.record?.id}
							<Icon icon="ph:at" />
							{formatTimeAgo(new Date())}
						</div>
					</div>
				</div>
				<Drawer.Close
					class={cn(
						buttonVariants({ size: 'icon', variant: 'ghost' }),
						'ml-auto rounded-full [&_svg]:size-6'
					)}
				>
					<Icon icon="ph:x-circle" />
				</Drawer.Close>
			</Drawer.Title>
			<Drawer.Description>
				<div class="grid gap-2 text-muted-foreground">
					<div class="grid grid-cols-1 items-center gap-2">
						<div class="flex items-center space-x-1">
							<Icon icon="ph-gear" class="size-5" />
							<span>postgres-to-postgres.connectors.pod.local</span>
						</div>
						<div class="flex items-center space-x-1">
							<Icon icon="ph:map-pin" class="size-5" />
							<span class="pr-4">1.2.3.4</span>
							<Icon icon="ph:copy" class="size-5" />
							<span>1 Replica</span>
						</div>
					</div>
				</div>
			</Drawer.Description>
		</Drawer.Header>
		<Tabs.Root value="jobs" class="px-8">
			<Tabs.List class="grid w-full grid-cols-3">
				<Tabs.Trigger value="jobs">Jobs 執行紀錄</Tabs.Trigger>
				<Tabs.Trigger value="mertics">Mertics 時間統計</Tabs.Trigger>
				<Tabs.Trigger value="configurations">Configurations + 更動的歷史</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="jobs" class="p-2">
				<Accordion.Root type="single" value="0">
					<Accordion.Item value="0">
						<Accordion.Trigger>{formatTimeAgo(new Date())}</Accordion.Trigger>
						<Accordion.Content>Latest</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item>
						<Accordion.Trigger>{formatTimeAgo(new Date(Date.now() - 86400000))}</Accordion.Trigger>
						<Accordion.Content>AAA</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item>
						<Accordion.Trigger
							>{formatTimeAgo(new Date(Date.now() - 86400000 * 2))}</Accordion.Trigger
						>
						<Accordion.Content>AAA</Accordion.Content>
					</Accordion.Item>
					<Accordion.Item>
						<Accordion.Trigger
							>{formatTimeAgo(new Date(Date.now() - 86400000 * 7))}</Accordion.Trigger
						>
						<Accordion.Content>AAA</Accordion.Content>
					</Accordion.Item>
				</Accordion.Root>
			</Tabs.Content>
			<Tabs.Content value="mertics" class="p-2">mertics</Tabs.Content>
			<Tabs.Content value="configurations" class="p-2"
				>configurations

				<Table.Root>
					<Table.Caption>A list of your recent invoices.</Table.Caption>
					<Table.Header>
						<Table.Row>
							<Table.Head class="w-[100px]">Invoice</Table.Head>
							<Table.Head>Status</Table.Head>
							<Table.Head>Method</Table.Head>
							<Table.Head class="text-right">Amount</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						<Table.Row>
							<Table.Cell class="font-medium">INV001</Table.Cell>
							<Table.Cell>Paid</Table.Cell>
							<Table.Cell>Credit Card</Table.Cell>
							<Table.Cell class="text-right">$250.00</Table.Cell>
						</Table.Row>
					</Table.Body>
				</Table.Root>
			</Tabs.Content>
		</Tabs.Root>
	</Drawer.Content>
</Drawer.Root>

<Dialog.Root
	onOpenChange={(open) => {
		if (!open) {
			setTimeout(() => {
				items = items.map((item) => ({ ...item, active: false }));
			}, 100);
		}
	}}
>
	<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>Create</Dialog.Trigger>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>
				{#each items as item}
					{#if item.active}
						<div class="flex items-center px-2">
							<Icon icon="ph:plugs" class="size-8" />
							<div class="space-x-2 pl-2">{item.name}</div>
						</div>
					{/if}
				{/each}
				{#if items.filter((item) => item.active).length === 0}
					<div class="flex items-center justify-center">What do you need?</div>
				{/if}
			</Dialog.Title>
			<Dialog.Description class="flex justify-center pt-4">
				{#if items[0].active}
					<FabricCreateSource item={items[0].set} />
					<!-- {:else if items[1].active}
					<FabricCreateDestination />
				{:else if items[2].active}
					<FabricCreateConnection /> -->
				{:else}
					<FabricCreateOverview bind:items />
				{/if}
			</Dialog.Description>
		</Dialog.Header>
	</Dialog.Content>
</Dialog.Root>
