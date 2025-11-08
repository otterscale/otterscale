<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { ScrollArea } from '$lib/components/ui/scroll-area';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { siteConfig } from '$lib/config/site';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { type Notification, notifications } from '$lib/stores';

	let { open = $bindable(false) }: { open: boolean } = $props();

	async function unreadNotifition(id: string) {
		notifications.update((items) => items.map((n) => (n.id === id ? { ...n, read: false } : n)));
	}
	async function readNotifition(id: string) {
		notifications.update((items) => items.map((n) => (n.id === id ? { ...n, read: true } : n)));
	}
	async function unarchiveNotifition(id: string) {
		notifications.update((items) =>
			items.map((n) => (n.id === id ? { ...n, archived: false } : n))
		);
	}
	async function archiveNotifition(id: string) {
		notifications.update((items) => items.map((n) => (n.id === id ? { ...n, archived: true } : n)));
	}
	async function deleteNotifition(id: string) {
		notifications.update((items) => items.map((n) => (n.id === id ? { ...n, deleted: true } : n)));
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="rounded-l-lg p-6">
		<Tabs.Root value="unread" class="space-y-2">
			<h1 class="text-xl font-semibold text-foreground">{m.notifications()}</h1>
			<Tabs.List class="grid w-full grid-cols-3">
				<Tabs.Trigger value="unread">{m.unread()}</Tabs.Trigger>
				<Tabs.Trigger value="archived">{m.archived()}</Tabs.Trigger>
				<Tabs.Trigger value="all">{m.all()}</Tabs.Trigger>
			</Tabs.List>
			<Separator />
			<Tabs.Content value="unread">
				{@render list($notifications.filter((n) => !n.read && !n.archived && !n.deleted))}
			</Tabs.Content>
			<Tabs.Content value="archived">
				{@render list($notifications.filter((n) => n.archived && !n.deleted))}
			</Tabs.Content>
			<Tabs.Content value="all">
				{@render list($notifications.filter((n) => !n.deleted))}
			</Tabs.Content>
		</Tabs.Root>
	</Sheet.Content>
</Sheet.Root>

{#snippet list(notifications: Notification[])}
	<ScrollArea class="h-screen">
		<div class="flex flex-col gap-4 py-4 pt-0">
			{#each notifications as notification}
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger>
							<div
								class="flex flex-col items-start gap-2 rounded-lg border p-4 text-left transition-all hover:bg-accent"
							>
								<div class="flex w-full items-center gap-2">
									{#if !notification.read}
										<span class="flex h-2 w-2 rounded-full bg-blue-600"></span>
									{/if}
									<div class="text-sm font-semibold">{notification.from ?? siteConfig.title}</div>

									<div class="ml-auto text-xs [&_svg]:size-5">
										<DropdownMenu.Root>
											<DropdownMenu.Trigger>
												<Icon icon="ph:dots-three" />
											</DropdownMenu.Trigger>
											<DropdownMenu.Content side="left" align="start">
												<DropdownMenu.Group>
													{#if notification.read}
														<DropdownMenu.Item
															onclick={() => {
																unreadNotifition(notification.id);
															}}
														>
															<Icon icon="ph:envelope-simple-bold" />
															<span>{m.mark_as_unread()}</span>
														</DropdownMenu.Item>
													{:else}
														<DropdownMenu.Item
															onclick={() => {
																readNotifition(notification.id);
															}}
														>
															<Icon icon="ph:envelope-simple-open-bold" />
															<span>{m.mark_as_read()}</span>
														</DropdownMenu.Item>
													{/if}
													{#if notification.archived}
														<DropdownMenu.Item
															onclick={() => {
																unarchiveNotifition(notification.id);
															}}
														>
															<Icon icon="ph:box-arrow-up-bold" />
															<span>{m.mark_as_unarchived()}</span>
														</DropdownMenu.Item>
													{:else}
														<DropdownMenu.Item
															onclick={() => {
																archiveNotifition(notification.id);
															}}
														>
															<Icon icon="ph:box-arrow-down-bold" />
															<span>{m.mark_as_archived()}</span>
														</DropdownMenu.Item>
													{/if}
													<DropdownMenu.Separator />
													<DropdownMenu.Item
														onclick={() => {
															deleteNotifition(notification.id);
														}}
													>
														<Icon icon="ph:trash-bold" class="text-red-500" />
														<span class="text-red-500">{m.delete()}</span>
													</DropdownMenu.Item>
												</DropdownMenu.Group>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</div>
								</div>
								<div class="flex flex-col space-y-1">
									<span class="text-xs font-medium">{notification.title}</span>
									<span class="line-clamp-2 text-left text-xs text-muted-foreground">
										{notification.content.substring(0, 300)}
									</span>
								</div>

								<div class="text-xs">
									{formatTimeAgo(notification.created)}
								</div>
							</div>
						</Tooltip.Trigger>
						<Tooltip.Content side="left" class="w-[350px] break-words whitespace-pre-wrap">
							<span class="break-words whitespace-pre-wrap">{notification.content}</span>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			{/each}
		</div>
	</ScrollArea>
{/snippet}
