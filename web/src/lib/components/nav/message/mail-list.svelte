<script lang="ts">
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import {
		archiveMessage,
		deleteMessage,
		readMessage,
		unarchiveMessage,
		unreadMessage,
		type pbMessage
	} from '$lib/pb.js';
	import Icon from '@iconify/svelte';
	import { toast } from 'svelte-sonner';

	export let type: number;
	export let items: pbMessage[];

	const DIVISIONS = [
		{ amount: 60, name: 'seconds' },
		{ amount: 60, name: 'minutes' },
		{ amount: 24, name: 'hours' },
		{ amount: 7, name: 'days' },
		{ amount: 4.34524, name: 'weeks' },
		{ amount: 12, name: 'months' },
		{ amount: Number.POSITIVE_INFINITY, name: 'years' }
	] as const;

	const formatter = new Intl.RelativeTimeFormat(undefined, {
		numeric: 'auto'
	});

	function formatTimeAgo(date: Date) {
		let duration = (date.getTime() - new Date().getTime()) / 1000;

		for (let i = 0; i <= DIVISIONS.length; i++) {
			const division = DIVISIONS[i];
			if (Math.abs(duration) < division.amount) {
				return formatter.format(Math.round(duration), division.name);
			}
			duration /= division.amount;
		}
	}

	function ok(msg: pbMessage): boolean {
		if (type === 0) {
			return !msg.isRead && !msg.isArchived && !msg.isDeleted;
		} else if (type === 1) {
			return msg.isArchived && !msg.isDeleted;
		}
		return !msg.isDeleted;
	}
</script>

<ScrollArea class="h-screen">
	<div class="flex flex-col gap-2 py-4 pt-0">
		{#each items as item}
			{#if ok(item)}
				<div
					class="flex flex-col items-start gap-2 rounded-lg border p-3 text-left text-sm transition-all hover:bg-accent"
				>
					<div class="flex w-full flex-col gap-1">
						<div class="flex items-center">
							<div class="flex items-center gap-2">
								<div class="font-semibold">{item.senderId}</div>
								{#if !item.isRead}
									<span class="flex h-2 w-2 rounded-full bg-blue-600"></span>
								{/if}
							</div>
							<div class="ml-auto text-xs">
								<DropdownMenu.Root>
									<DropdownMenu.Trigger>
										<Icon icon="ph:dots-three-vertical-bold" class="h-5 w-5" />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content>
										<DropdownMenu.Group>
											{#if item.isRead}
												<DropdownMenu.Item
													class="space-x-1"
													on:click={async () => {
														await unreadMessage(item.id);
														item.isRead = false;
														items = items;
														toast.success('Message marked as unread.');
													}}
												>
													<Icon icon="ph:arrow-counter-clockwise" class="h-5 w-5" />
													<span>Unread</span>
												</DropdownMenu.Item>
											{:else}
												<DropdownMenu.Item
													class="space-x-1"
													on:click={async () => {
														await readMessage(item.id);
														item.isRead = true;
														items = items;
														toast.success('Message marked as read.');
													}}
												>
													<Icon icon="ph:check-circle" class="h-5 w-5" />
													<span>Read</span>
												</DropdownMenu.Item>
											{/if}
											{#if item.isArchived}
												<DropdownMenu.Item
													class="space-x-1"
													on:click={async () => {
														await unarchiveMessage(item.id);
														item.isArchived = false;
														items = items;
														toast.success('Message removed from archive.');
													}}
												>
													<Icon icon="ph:box-arrow-up" class="h-5 w-5" />
													<span>Unarchive</span>
												</DropdownMenu.Item>
											{:else}
												<DropdownMenu.Item
													class="space-x-1"
													on:click={async () => {
														await archiveMessage(item.id);
														item.isArchived = true;
														items = items;
														toast.success('Message moved to archive.');
													}}
												>
													<Icon icon="ph:box-arrow-down" class="h-5 w-5" />
													<span>Archive</span>
												</DropdownMenu.Item>
											{/if}
											<DropdownMenu.Item
												class="space-x-1"
												on:click={async () => {
													await deleteMessage(item.id);
													item.isDeleted = true;
													items = items;
													toast.success('Message deleted.');
												}}
											>
												<Icon icon="ph:trash" class="h-5 w-5" />
												<span>Delete</span>
											</DropdownMenu.Item>
										</DropdownMenu.Group>
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							</div>
						</div>
						<div class="text-xs font-medium">{item.title}</div>
					</div>
					<div class="line-clamp-2 text-xs text-muted-foreground">
						{item.content.substring(0, 300)}
					</div>
					<div class="text-xs">
						{formatTimeAgo(new Date(item.created))}
					</div>
				</div>
			{/if}
		{/each}
	</div>
</ScrollArea>
