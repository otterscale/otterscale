<script lang="ts">
	import Icon from '@iconify/svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import * as Accordion from '$lib/components/ui/accordion';
	import * as Drawer from '$lib/components/ui/drawer';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';

	import { formatTimeAgo } from '$lib/formatter';
	import { cn } from '$lib/utils';
	import type { Node } from '@xyflow/svelte';
	import { connectorIcon } from '$lib/connector';

	let {
		open = $bindable(),
		node
	}: {
		open: boolean;
		node: Node;
	} = $props();
</script>

<Drawer.Root direction="right" bind:open>
	<Drawer.Trigger class={buttonVariants({ variant: 'outline' })}>Open Drawer</Drawer.Trigger>
	<Drawer.Content class="absolute inset-x-auto inset-y-0 right-0 w-3/5 space-y-2 rounded-tr-none">
		<Drawer.Header class="px-8 pt-0">
			<Drawer.Title class="flex items-center">
				<div class="flex items-center space-x-2">
					<Icon icon={connectorIcon(node.data.type as string)} class="size-8" />
					<div class="flex-col p-2">
						{node.data.name}
						<div class="flex text-sm text-muted-foreground">
							{node.data.type}
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
				<div class="flex-col space-y-1">
					<Badge variant="outline" class="text-muted-foreground hover:scale-105">
						<Icon icon="ph:folder" class="size-5 sm:flex" />
						<span class="pl-1 text-sm">{node.data.kind}</span>
					</Badge>
					<div>
						<Badge variant="outline" class="text-muted-foreground hover:scale-105">
							<Icon icon="ph:gear" class="size-5 sm:flex" />
							<span class="pl-1 text-sm">postgres-to-postgres.connectors.pod.local</span>
						</Badge>
					</div>
					<div>
						<Badge variant="outline" class="text-muted-foreground hover:scale-105">
							<Icon icon="ph:floppy-disk-back" class="size-5 sm:flex" />
							<span class="pl-1 text-sm">{node.data.image}</span>
						</Badge>
						<Badge variant="outline" class="text-muted-foreground hover:scale-105">
							<Icon icon="ph:map-pin" class="size-5 sm:flex" />
							<span class="pl-1 text-sm">1.2.3.4</span>
						</Badge>
						<Badge variant="outline" class="text-muted-foreground hover:scale-105">
							<Icon icon="ph:copy" class="size-5 sm:flex" />
							<span class="pl-1 text-sm">1 Replica</span>
						</Badge>
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
