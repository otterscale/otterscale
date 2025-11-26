<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import SendIcon from '@lucide/svelte/icons/send';

	import * as Chat from '$lib/components/custom/chat';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';

	import type { Choice, Message } from './types.d';

	const userIdentifier = 'userIdentifier';
	const receiverIdentifier = 'receiverIdentifier';
</script>

<script lang="ts">
	let {
		serviceURL,
		modelName,
		namespace
	}: { serviceURL: string; modelName: string; namespace: string } = $props();

	let open = $state(false);
	function close() {
		open = false;
	}

	let message = $state('');
	let messages: Message[] = $state([]);

	let isSubmitting = $state(false);
	async function onsubmit(model: string, prompt: string) {
		if (isSubmitting) return;
		isSubmitting = true;

		try {
			const response = await fetch(`${serviceURL}/v1/completions`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ model, prompt, max_tokens: 40 })
			});

			if (!response.ok) {
				console.error('Request failed:', response.status);
				return;
			}

			const data = await response.json();
			const text = data.choices?.[0]?.text || '';

			messages.push({
				senderId: receiverIdentifier,
				message: text,
				sentAt: new Date().toLocaleTimeString('en-US', {
					hour: 'numeric',
					minute: '2-digit'
				})
			});
		} catch (error) {
			console.error('Failed to fetch completion:', error);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		<Icon icon="ph:robot" />
	</AlertDialog.Trigger>
	<AlertDialog.Content class="max-w-[77vw] min-w-[77vw] p-1">
		<div class="rounded-lg border-border">
			<div
				class="flex h-[70px] place-items-center justify-between rounded-t-lg border-b bg-background p-2"
			>
				<div class="flex place-items-center gap-2">
					<div class="relative flex size-8 shrink-0 overflow-hidden rounded-full border bg-muted">
						<Icon
							icon="ph:robot-bold"
							class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
						/>
					</div>
					<div class="flex flex-col">
						<span class="text-sm font-medium">{modelName}</span>
						<span class="text-xs">{namespace}</span>
					</div>
				</div>
				<div class="flex place-items-center gap-1">
					{#if messages.length > 0}
						<Button
							disabled={isSubmitting}
							variant="ghost"
							size="icon"
							class="rounded-full"
							onclick={() => {
								messages = [] as Message[];
							}}
						>
							<Icon icon="ph:arrows-clockwise-bold" />
						</Button>
					{/if}
					<Button
						variant="ghost"
						size="icon"
						class="rounded-full"
						onclick={() => {
							close();
						}}
					>
						<Icon icon="ph:x-bold" />
					</Button>
				</div>
			</div>
			{#if messages.length === 0}
				<div
					class="relative mx-auto mt-0 flex h-[calc(77vh-130px)] w-full max-w-4xl flex-col items-center justify-center px-4 text-center sm:px-6 lg:px-8"
				>
					<Icon
						icon="ph:sparkle"
						class="absolute -z-10 h-[320px] w-[320px] rotate-45 transform animate-pulse text-muted-foreground opacity-10 blur-sm"
					/>
					<div class="z-10">
						<h1 class="text-3xl font-bold text-primary sm:text-4xl">Model Testing Sandbox</h1>
						<p class="mt-3 text-muted-foreground">
							Type a prompt below and press Send to test the model.
						</p>
					</div>
				</div>
			{:else}
				<Chat.List class="h-[calc(77vh-130px)] max-h-[calc(77vh-130px)]">
					{#each messages as message (message)}
						<Chat.Bubble variant={message.senderId === userIdentifier ? 'sent' : 'received'}>
							<div
								class="relative order-1 flex size-8 shrink-0 overflow-hidden rounded-full border group-data-[variant='sent']/chat-bubble:order-2"
							>
								<Icon
									icon={message.senderId === userIdentifier ? 'ph:user' : 'ph:robot'}
									class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
								/>
							</div>
							<Chat.BubbleMessage class="flex flex-col gap-1">
								<p>{message.message}</p>
								<div class="w-full text-xs group-data-[variant='sent']/chat-bubble:text-end">
									{message.sentAt}
								</div>
							</Chat.BubbleMessage>
						</Chat.Bubble>
					{/each}
					{#if messages[messages.length - 1].senderId === userIdentifier}
						<Chat.Bubble variant="received">
							<div class="relative flex size-8 shrink-0 overflow-hidden rounded-full border">
								<Icon
									icon="ph:robot"
									class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
								/>
							</div>
							<Chat.BubbleMessage typing />
						</Chat.Bubble>
					{/if}
				</Chat.List>
			{/if}
			<form
				onsubmit={(e) => {
					e.preventDefault();
					messages.push({
						message,
						senderId: userIdentifier,
						sentAt: new Date().toLocaleTimeString('en-US', {
							hour: 'numeric',
							minute: '2-digit'
						})
					});
					message = '';
				}}
				class="flex h-[70px] w-full place-items-center gap-2 p-2"
			>
				<Input
					bind:value={message}
					class="rounded-full disabled:cursor-default"
					placeholder="Type a message..."
				/>
				<Button
					type="submit"
					variant="default"
					size="icon"
					class="shrink-0 rounded-full"
					disabled={message === ''}
					onclick={() => {
						onsubmit(modelName, message);
					}}
				>
					<SendIcon />
				</Button>
			</form>
		</div>
	</AlertDialog.Content>
</AlertDialog.Root>
