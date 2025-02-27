<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Drawer from '$lib/components/ui/drawer';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Accordion from '$lib/components/ui/accordion';
	import { formatTimeAgo } from '$lib/formatter';
	import pb from '$lib/pb';

	let open = false;
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
				<Accordion.Root type="single">
					<Accordion.Item value="item-1">
						<Accordion.Trigger>Is it accessible?</Accordion.Trigger>
						<Accordion.Content>Yes. It adheres to the WAI-ARIA design pattern.</Accordion.Content>
					</Accordion.Item>
				</Accordion.Root>
			</Tabs.Content>
			<Tabs.Content value="mertics" class="p-2">mertics</Tabs.Content>
			<Tabs.Content value="configurations" class="p-2">configurations</Tabs.Content>
		</Tabs.Root>
	</Drawer.Content>
</Drawer.Root>
