<script lang="ts">
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';

	import { buttonVariants } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as m from '$lib/paraglide/messages.js';
	import { listMessages, type pbMessage } from '$lib/pb';

	import MailList from './mail-list.svelte';
	import { cn } from '$lib/utils';

	let msgs: pbMessage[] = [];
	onMount(async () => {
		msgs = await listMessages();
	});
</script>

<Sheet.Root>
	<Sheet.Trigger
		class={cn(buttonVariants({ variant: 'outline', size: 'icon' }), 'bg-header [&_svg]:size-5')}
	>
		{#if msgs.filter((msg) => !msg.isRead).length > 0}
			<Icon icon="ph:notification-fill" />
		{:else}
			<Icon icon="ph:notification" />
		{/if}
	</Sheet.Trigger>
	<Sheet.Content class="rounded-l-[10px]">
		<Tabs.Root value="unread">
			<h1 class="text-xl font-bold">{m.inbox()}</h1>
			<div class="flex items-center py-4">
				<Tabs.List class="grid w-full grid-cols-3">
					<Tabs.Trigger value="unread">{m.inbox_unread()}</Tabs.Trigger>
					<Tabs.Trigger value="archived">{m.inbox_archived()}</Tabs.Trigger>
					<Tabs.Trigger value="all">{m.inbox_all()}</Tabs.Trigger>
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
