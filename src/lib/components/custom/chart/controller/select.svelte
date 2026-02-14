<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';

	let {
		timeRange = $bindable(),
		options = [
			{ value: '90d', label: 'Last 3 months' },
			{ value: '30d', label: 'Last 30 days' },
			{ value: '7d', label: 'Last 7 days' }
		]
	}: {
		timeRange: string;
		options?: Array<{ value: string; label: string }>;
	} = $props();

	const selectedLabel = $derived.by(() => {
		const option = options.find((opt) => opt.value === timeRange);
		return option?.label || options[0]?.label || '';
	});
</script>

<div class="flex px-6">
	<Select.Root type="single" bind:value={timeRange}>
		<Select.Trigger class="w-[160px] rounded-lg sm:ml-auto" aria-label="Select a value">
			{selectedLabel}
		</Select.Trigger>
		<Select.Content class="rounded-xl">
			{#each options as option}
				<Select.Item value={option.value} class="rounded-lg">{option.label}</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
</div>
