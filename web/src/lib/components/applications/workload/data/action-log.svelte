<script lang="ts">
	import { timestampFromDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onDestroy } from 'svelte';
	import { fly } from 'svelte/transition';

	import type { Application_Pod } from '$lib/api/application/v1/application_pb';
	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Select from '$lib/components/ui/select';
	import * as Sheet from '$lib/components/ui/sheet';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';

	let { namespace, pod }: { namespace: string; pod: Application_Pod } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	let open = $state(false);
	let messages = $state<{ message: string; phase: string }[]>([]);
	let shouldStop = $state(false);
	let terminal = $state<HTMLDivElement>();

	// Time options for log filtering
	type TimeOption = {
		label: string;
		value: string;
		minutes: number;
	};

	const timeOptions: TimeOption[] = [
		{ label: 'All', value: 'all', minutes: 0 },
		{ label: '5m', value: '5m', minutes: 5 },
		{ label: '15m', value: '15m', minutes: 15 },
		{ label: '30m', value: '30m', minutes: 30 },
		{ label: '1h', value: '1h', minutes: 60 },
		{ label: '3h', value: '3h', minutes: 180 },
		{ label: '6h', value: '6h', minutes: 360 },
		{ label: '12h', value: '12h', minutes: 720 },
	];

	let selectedTimeValue = $state<string>(timeOptions[0].value); // For Select component
	let selectedTime = $derived(timeOptions.find((option) => option.value === selectedTimeValue)); // Derived from selectedTimeValue

	async function watchLogs() {
		try {
			messages = [];
			shouldStop = false;

			const since =
				selectedTime && selectedTime.minutes > 0
					? timestampFromDate(new Date(Date.now() - selectedTime.minutes * 60 * 1000))
					: undefined;

			const stream = applicationClient.watchLogs({
				scope: $currentKubernetes?.scope || '',
				facility: $currentKubernetes?.name || '',
				namespace,
				podName: pod.name,
				containerName: '',
				since,
			});
			for await (const response of stream) {
				if (shouldStop) break;

				messages = [...messages, { message: response.log, phase: 'LOG' }];
				scrollToBottom();

				// Add delay for better UX
				await new Promise((resolve) => setTimeout(resolve, 100));
			}
		} catch (error) {
			// Ignore abort-related errors
			const isAbortError =
				error instanceof Error &&
				(error.name === 'AbortError' ||
					error.message.includes('aborted') ||
					error.message.includes('canceled'));

			if (!isAbortError) {
				console.error('Error watching logs:', error);
			}
		}
	}

	function scrollToBottom() {
		if (terminal) {
			terminal.scrollTop = terminal.scrollHeight;
		}
	}

	function stopWatching() {
		shouldStop = true;
		messages = [];
	}

	// Stop watching when sheet closes
	$effect(() => {
		if (!open && shouldStop === false) {
			stopWatching();
		}
	});

	$effect(() => {
		if (open && selectedTime) {
			stopWatching();
			setTimeout(watchLogs, 100);
		}
	});

	onDestroy(() => {
		stopWatching();
	});
</script>

<Sheet.Root bind:open>
	<Sheet.Trigger onclick={watchLogs}>
		<span class="flex items-center gap-1">
			<Icon icon="ph:file-text" />
			{m.log()}
		</span>
	</Sheet.Trigger>

	<Sheet.Content class="rounded-l-lg border-none bg-transparent sm:max-w-9/10">
		{#if open}
			<div
				class="border-border bg-secondary text-card-foreground dark size-full flex-col rounded-l-lg border font-mono text-sm shadow-sm"
			>
				<!-- Header with time selector -->
				<div class="border-border flex items-center justify-between border-b p-4">
					<h3 class="text-lg font-semibold">{m.log()}</h3>
					<div class="mr-8 flex items-center gap-2">
						<span class="text-muted-foreground text-sm">Since:</span>
						<Select.Root type="single" bind:value={selectedTimeValue}>
							<Select.Trigger class="w-20">
								{#if selectedTime}
									{selectedTime.label}
								{:else}
									<span>Select time</span>
								{/if}
							</Select.Trigger>
							<Select.Content>
								{#each timeOptions as option}
									<Select.Item value={option.value}>
										{option.label}
									</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				</div>

				<div bind:this={terminal} class="h-full flex-col overflow-auto p-4">
					{#each messages as msg, i}
						{@const isLastMessage = messages.length === i + 1}
						{@const textClass = isLastMessage ? '' : 'text-green-500'}

						{#if msg.message}
							<span class="block" transition:fly={{ y: -5, duration: 500 }}>
								<span class="flex space-x-1 {textClass}">
									<span>[{msg.phase}] {msg.message}</span>
								</span>
							</span>
						{/if}
					{/each}
				</div>
			</div>
		{/if}
	</Sheet.Content>
</Sheet.Root>
