<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import SendIcon from '@lucide/svelte/icons/send';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Chat from '$lib/components/custom/chat';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import { Label } from '$lib/components/ui/label';

	import type { Message } from './types.d';

	const userIdentifier = 'user';
	const receiverIdentifier = 'receiver';
	const defaults = {
		temperature: 0.1,
		max_tokens: 128
	};
</script>

<script lang="ts">
	let { serviceUri, model }: { serviceUri: string; model: Model } = $props();

	// Parameters
	let temperature = $state(defaults.temperature);
	let max_tokens = $state(defaults.max_tokens);

	// Messages
	let messages: Message[] = $state([]);
	const isNewChat = $derived(messages.length === 0);

	// Compltion
	let hasError = $state(false);
	let userMessage = $state('');
	let modelMessage = $state('');
	let isModelLoaded = $state(false);
	async function onsubmit() {
		isModelLoaded = false;
		messages.push({
			message: userMessage,
			senderId: userIdentifier,
			sentAt: new Date().toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit'
			})
		});

		const response = await fetch('/api/completion', {
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
		});

		const body = await response.json();

		if (!response.ok) {
			hasError = true;
			throw new Error('Failed to get response from model:', body);
		}

		modelMessage = body.choices.map((choice: { text: string }) => choice.text).join('');

		isModelLoaded = true;
		messages.push({
			message: modelMessage,
			senderId: receiverIdentifier,
			sentAt: new Date().toLocaleTimeString('en-US', {
				hour: 'numeric',
				minute: '2-digit'
			})
		});

		userMessage = '';
		modelMessage = '';
	}

	function reset() {
		max_tokens = defaults.max_tokens;
		temperature = defaults.temperature;
		userMessage = '';
		modelMessage = '';
		messages = [] as Message[];
		isModelLoaded = false;
		hasError = false;
	}
</script>

<Dialog.Root
	onOpenChangeComplete={() => {
		reset();
	}}
>
	<Dialog.Trigger class={buttonVariants({ variant: 'ghost' })}>
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
							<Dialog.Title>Controllers</Dialog.Title>
						</Dialog.Header>
						<div class="flex flex-col gap-4 rounded-lg border-border bg-background">
							<div class="flex flex-col gap-2">
								<div class="flex justify-between gap-4">
									<Label for="temperature" class="text-sm font-medium">Temperature</Label>
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
									<Label for="maximum_token_length" class="text-sm font-medium">Max Tokens</Label>
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
							reset();
						}}
					>
						<Icon icon="ph:arrows-clockwise-bold" class="size-5" />
					</Button>
				{/if}
			</div>
		</div>
		<div class="relative h-[calc(77vh-200px)] overflow-y-auto">
			{#if isNewChat}
				<!-- Background -->
				<div
					class="relative mx-auto mt-0 flex h-[50vh] w-full max-w-4xl flex-col items-center justify-center px-4 text-center sm:px-6 lg:px-8"
				>
					<Icon
						icon="ph:sparkle"
						class="absolute -z-10 h-[500px] w-[500px] rotate-45 transform animate-pulse text-muted-foreground opacity-10 blur-sm"
					/>
					<div class="z-10">
						<h1 class="text-3xl font-bold text-primary sm:text-4xl">Model Testing</h1>
						<p class="mt-3 text-muted-foreground">
							Type a prompt below and press Send to test the model.
						</p>
					</div>
				</div>
			{:else}
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
								<p class="text-destructive">An error occurred while fetching the response.</p>
							</Chat.BubbleMessage>
						</Chat.Bubble>
					{:else if !isModelLoaded && modelMessage === ''}
						<Chat.Bubble variant="received">
							<div class="relative flex size-8 shrink-0 overflow-hidden rounded-full border">
								<Icon
									icon="ph:robot"
									class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
								/>
							</div>
							<Chat.BubbleMessage typing />
						</Chat.Bubble>
					{:else if !isModelLoaded && modelMessage !== ''}
						<Chat.Bubble variant="received">
							<div
								class="relative order-1 flex size-8 shrink-0 overflow-hidden rounded-full border group-data-[variant='sent']/chat-bubble:order-2"
							>
								<Icon
									icon="ph:robot"
									class="absolute top-1/2 left-1/2 size-5 -translate-x-1/2 -translate-y-1/2"
								/>
							</div>
							<Chat.BubbleMessage class="flex max-w-96 flex-col gap-1 break-all">
								<p>{modelMessage}</p>
							</Chat.BubbleMessage>
						</Chat.Bubble>
					{/if}
				</Chat.List>
			{/if}
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
					class="rounded-lg text-destructive/70 hover:text-primary"
					onclick={() => {
						userMessage = 'Are you alive?';
						onsubmit();
					}}
				>
					<Icon icon="ph:heartbeat" />
					<p class="text-xs">Health Check</p>
				</InputGroup.Button>
				<InputGroup.Button
					variant="default"
					class="ml-auto rounded-full"
					size="icon-sm"
					disabled={userMessage === ''}
					onclick={() => {
						onsubmit();
					}}
				>
					<SendIcon />
				</InputGroup.Button>
			</InputGroup.Addon>
		</InputGroup.Root>
	</Dialog.Content>
</Dialog.Root>
