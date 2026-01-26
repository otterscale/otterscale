<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import SendIcon from '@lucide/svelte/icons/send';
	import { Dialog as DialogPrimitive } from 'bits-ui';
	import { toast } from 'svelte-sonner';

	import { resolve } from '$app/paths';
	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Chat from '$lib/components/custom/chat';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';

	import type { Message } from './types.d';

	const userIdentifier = 'user';
	const receiverIdentifier = 'receiver';
	const defaults = {
		temperature: 0.1,
		max_tokens: 128
	};
</script>

<script lang="ts">
	let {
		serviceUri,
		model,
		scope,
		...restProps
	}: DialogPrimitive.TriggerProps & { serviceUri: string; model: Model; scope: string } = $props();

	let conversation = $state<HTMLDivElement>();

	// Parameters
	let temperature = $state(defaults.temperature);
	let max_tokens = $state(defaults.max_tokens);

	// Messages
	let messages: Message[] = $state([
		{
			message: m.model_testing_init(),
			senderId: receiverIdentifier,
			sentAt: new Date().toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit'
			})
		}
	]);
	const isNewChat = $derived(messages.length === 0);
	const latestModelResponseIndex = $derived(
		messages.findLastIndex((message) => message.senderId === receiverIdentifier)
	);

	// Completion
	let isModelLoaded = $state(true);
	let hasError = $state(false);
	let userMessage = $state('');
	let currentAbortController = $state({} as AbortController);
	async function onsubmit() {
		isModelLoaded = false;

		messages = [
			...messages,
			{
				message: userMessage,
				senderId: userIdentifier,
				sentAt: new Date().toLocaleTimeString('en-US', {
					hour: 'numeric',
					minute: '2-digit'
				})
			}
		];

		const response = await fetch(
			resolve('/(auth)/scope/[scope]/models/api/completion', { scope }),
			{
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					serviceUri: serviceUri,
					modelName: model.name,
					modelIdentifier: model.id,
					prompt: userMessage,
					max_tokens: max_tokens,
					temperature: temperature
				})
			}
		);

		if (!response.ok) {
			hasError = true;
			const errorMessage = await response.text();
			toast.error('Failed to get response from model', {
				description: errorMessage,
				duration: Number.POSITIVE_INFINITY,
				closeButton: true
			});
		}

		if (response.body) {
			const reader = response.body.getReader();
			const decoder = new TextDecoder();
			let buffer = '';

			const abortController = new AbortController();
			currentAbortController = abortController;

			messages = [
				...messages,
				{
					message: '',
					senderId: receiverIdentifier,
					sentAt: new Date().toLocaleTimeString('en-US', {
						hour: 'numeric',
						minute: '2-digit'
					})
				}
			];
			try {
				while (true) {
					if (abortController.signal.aborted) break;

					const { done, value } = await reader.read();
					if (done) break;

					buffer += decoder.decode(value, { stream: true });
					const lines = buffer.split('\n');
					buffer = lines.pop() || '';

					for (const line of lines) {
						if (line.startsWith('data: ')) {
							const jsonString = line.slice(6);
							if (jsonString === '[DONE]') {
								continue;
							}
							try {
								const data = JSON.parse(jsonString);
								if (data.choices?.[0]?.text) {
									messages[latestModelResponseIndex].message += data.choices[0].text;
								}
							} catch (error) {
								console.error(error);
							}

							scrollToBottom();
						}
					}
				}
			} catch (error) {
				if (error instanceof Error && error.name !== 'AbortError') {
					console.error(error);
				}
			} finally {
				reader.releaseLock();
			}
		}

		userMessage = '';
		isModelLoaded = true;
	}

	function init() {
		max_tokens = defaults.max_tokens;
		temperature = defaults.temperature;
		userMessage = '';
		messages = [
			{
				message: `Welcome! I'm your AI assistant. How can I help you today?`,
				senderId: receiverIdentifier,
				sentAt: new Date().toLocaleTimeString('en-US', {
					hour: 'numeric',
					minute: '2-digit'
				})
			}
		] as Message[];
		isModelLoaded = true;
		hasError = false;
	}

	function scrollToBottom() {
		if (conversation) {
			conversation.scrollTop = conversation.scrollHeight;
		}
	}
</script>

<Dialog.Root onOpenChangeComplete={init}>
	<Dialog.Trigger class={buttonVariants({ variant: 'ghost' })} {...restProps}>
		<Icon icon="ph:robot" />
	</Dialog.Trigger>
	<Dialog.Content
		class="flex h-[77vh] max-w-[77vw] min-w-[50vw] flex-col justify-between overflow-auto"
	>
		<div
			class="flex h-12 place-items-center justify-between rounded-t-lg border-b bg-background p-2"
		>
			<!-- Model Informations -->
			<div class="flex place-items-center gap-2">
				<div class="relative flex size-8 shrink-0 overflow-hidden rounded-full border bg-muted">
					<Icon
						icon="ph:robot-bold"
						class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
					/>
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-medium">{model.name}</span>
					<span class="text-xs text-muted-foreground">{model.id}</span>
				</div>
			</div>
			<!-- Controllers -->
			<div class="flex place-items-center gap-1">
				<Dialog.Root>
					<Dialog.Trigger
						class={buttonVariants({ variant: 'ghost', size: 'icon', class: 'rounded-full' })}
					>
						<Icon icon="ph:faders-horizontal-bold" class="size-5" />
					</Dialog.Trigger>
					<Dialog.Content>
						<Dialog.Header>
							<Dialog.Title>{m.parameters()}</Dialog.Title>
						</Dialog.Header>
						<div class="flex flex-col gap-4 rounded-lg border-border bg-background">
							<div class="flex flex-col gap-2">
								<div class="flex justify-between gap-4">
									<Label for="temperature" class="text-sm font-medium">{m.temperature()}</Label>
									<p class="text-muted-foreground">{temperature}</p>
								</div>
								<Input
									type="number"
									min="0"
									max="1"
									step="0.01"
									bind:value={temperature}
									class="w-full"
								/>
							</div>
							<div class="flex flex-col gap-2">
								<div class="flex justify-between gap-4">
									<Label for="maximum_token_length" class="text-sm font-medium"
										>{m.max_tokens()}</Label
									>
									<p class="text-muted-foreground">{max_tokens}</p>
								</div>
								<Input
									type="number"
									min="0"
									max="1024"
									step="1"
									bind:value={max_tokens}
									class="w-full"
								/>
							</div>
						</div>
					</Dialog.Content>
				</Dialog.Root>
				{#if !isNewChat}
					<Button
						disabled={!isModelLoaded && !hasError}
						variant="ghost"
						size="icon"
						class="rounded-full"
						onclick={() => {
							init();
						}}
					>
						<Icon icon="ph:arrows-clockwise-bold" class="size-5" />
					</Button>
				{/if}
			</div>
		</div>
		<div bind:this={conversation} class="relative h-[calc(77vh-200px)] overflow-y-auto">
			<!-- Conversations -->
			<Chat.List>
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
						<Chat.BubbleMessage class="flex max-w-96 flex-col gap-1 break-all">
							<p>{message.message}</p>
							<div class="w-full text-xs group-data-[variant='sent']/chat-bubble:text-end">
								{message.sentAt}
							</div>
						</Chat.BubbleMessage>
					</Chat.Bubble>
				{/each}
				{#if hasError}
					<Chat.Bubble variant="received">
						<div
							class="relative order-1 flex size-8 shrink-0 overflow-hidden rounded-full border group-data-[variant='sent']/chat-bubble:order-2"
						>
							<Icon
								icon="ph:robot"
								class="absolute  top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
							/>
						</div>
						<Chat.BubbleMessage class="flex max-w-96 flex-col gap-1 break-all">
							<p class="text-destructive">
								{m.model_response_error()}
							</p>
						</Chat.BubbleMessage>
					</Chat.Bubble>
				{:else if !isModelLoaded && messages[latestModelResponseIndex].message === ''}
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
		</div>
		<!-- Inputs -->
		<InputGroup.Root class="flex h-24 bg-muted">
			<InputGroup.Textarea
				placeholder="Start chat with {model.id} {model.name}..."
				bind:value={userMessage}
			/>
			<InputGroup.Addon align="block-end" class="flex h-12 items-center gap-2">
				<InputGroup.Button
					variant="outline"
					class="text-muted-foreground hover:text-primary"
					disabled={hasError}
					onclick={() => {
						userMessage = m.health_check_statment();
						onsubmit();
					}}
				>
					<Icon icon="ph:heartbeat" />
					<p class="text-xs">{m.health_check()}</p>
				</InputGroup.Button>

				{#if isModelLoaded}
					<Button
						variant="default"
						class="ml-auto rounded-full"
						size="icon-sm"
						disabled={hasError || userMessage === ''}
						onclick={() => {
							onsubmit();
						}}
					>
						<SendIcon />
					</Button>
				{:else}
					<Button
						variant="default"
						class="ml-auto rounded-full"
						size="icon-sm"
						disabled={hasError}
						onclick={() => {
							currentAbortController.abort();
						}}
					>
						<Icon icon="ph:stop" />
					</Button>
				{/if}
			</InputGroup.Addon>
		</InputGroup.Root>
	</Dialog.Content>
</Dialog.Root>
