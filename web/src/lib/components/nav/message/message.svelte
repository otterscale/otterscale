<script lang="ts">
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';

	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { listMessages, type pbMessage } from '$lib/pb';
	import MailList from './mail-list.svelte';

	let msgs: pbMessage[] = [];
	onMount(async () => {
		msgs = await listMessages();
	});
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<Tooltip.Root>
			<Tooltip.Trigger asChild>
				<Button variant="outline" size="icon" class="bg-header">
					{#if msgs.filter((msg) => !msg.isRead).length > 0}
						<Icon icon="ph:notification-fill" class="h-5 w-5" />
					{:else}
						<Icon icon="ph:notification" class="h-5 w-5" />
					{/if}
				</Button>
			</Tooltip.Trigger>
			<Tooltip.Content>
				<p>Inbox</p>
			</Tooltip.Content>
		</Tooltip.Root>
	</Sheet.Trigger>
	<Sheet.Content>
		<Tabs.Root value="unread">
			<h1 class="text-xl font-bold">Inbox</h1>
			<div class="flex items-center py-4">
				<Tabs.List class="grid w-full grid-cols-3">
					<Tabs.Trigger value="unread">Unread</Tabs.Trigger>
					<Tabs.Trigger value="archived">Archived</Tabs.Trigger>
					<Tabs.Trigger value="all">All</Tabs.Trigger>
				</Tabs.List>
			</div>
			<Separator />
			<Tabs.Content value="unread" class="py-2">
				<MailList bind:items={msgs} type={0} />
			</Tabs.Content>
			<Tabs.Content value="archived" class="py-2">
				<MailList bind:items={msgs} type={1} />
			</Tabs.Content>
			<Tabs.Content value="all" class="py-2">
				<MailList bind:items={msgs} type={2} />
			</Tabs.Content>
		</Tabs.Root>
	</Sheet.Content>
</Sheet.Root>
