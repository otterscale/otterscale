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
	import { formatTimeAgo } from '$lib/formatter';
	import { siteConfig } from '$lib/config/site';

	export let type: number;
	export let items: pbMessage[];

	function ok(msg: pbMessage): boolean {
		if (type === 0) {
			return !msg.read && !msg.archived && !msg.deleted;
		} else if (type === 1) {
			return msg.archived && !msg.deleted;
		}
		return !msg.deleted;
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
								<div class="font-semibold">{item.from ?? siteConfig.title}</div>
								{#if !item.read}
									<span class="flex h-2 w-2 rounded-full bg-blue-600"></span>
								{/if}
							</div>
							<div class="ml-auto text-xs [&_svg]:size-5">
								<DropdownMenu.Root>
									<DropdownMenu.Trigger>
										<Icon icon="ph:dots-three-vertical-bold" />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content>
										<DropdownMenu.Group>
											{#if item.read}
												<DropdownMenu.Item
													class="space-x-1"
													onclick={async () => {
														await unreadMessage(item.id);
														item.read = false;
														items = items;
														toast.success('Message marked as unread.');
													}}
												>
													<Icon icon="ph:arrow-counter-clockwise" />
													<span>Unread</span>
												</DropdownMenu.Item>
											{:else}
												<DropdownMenu.Item
													class="space-x-1"
													onclick={async () => {
														await readMessage(item.id);
														item.read = true;
														items = items;
														toast.success('Message marked as read.');
													}}
												>
													<Icon icon="ph:check-circle" />
													<span>Read</span>
												</DropdownMenu.Item>
											{/if}
											{#if item.archived}
												<DropdownMenu.Item
													class="space-x-1"
													onclick={async () => {
														await unarchiveMessage(item.id);
														item.archived = false;
														items = items;
														toast.success('Message removed from archive.');
													}}
												>
													<Icon icon="ph:box-arrow-up" />
													<span>Unarchive</span>
												</DropdownMenu.Item>
											{:else}
												<DropdownMenu.Item
													class="space-x-1"
													onclick={async () => {
														await archiveMessage(item.id);
														item.archived = true;
														items = items;
														toast.success('Message moved to archive.');
													}}
												>
													<Icon icon="ph:box-arrow-down" />
													<span>Archive</span>
												</DropdownMenu.Item>
											{/if}
											<DropdownMenu.Item
												class="space-x-1"
												onclick={async () => {
													await deleteMessage(item.id);
													item.deleted = true;
													items = items;
													toast.success('Message deleted.');
												}}
											>
												<Icon icon="ph:trash" />
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
