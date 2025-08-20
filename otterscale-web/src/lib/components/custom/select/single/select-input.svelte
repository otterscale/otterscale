<script lang="ts">
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';
	import { getContext, type Snippet } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { OptionType } from './types';
	import type { OptionManager } from './utils.svelte';

	const optionManager: OptionManager = getContext('OptionManager');
	let {
		ref = $bindable(null),
		class: className,
		value = $bindable(''),
		addition,
		...restProps
	}: CommandPrimitive.InputProps & {
		addition?: Snippet<
			[
				{
					manager: OptionManager;
					accessor: {
						input: string;
					};
				}
			]
		>;
	} = $props();

	let options: Writable<OptionType[]> = getContext('options');
	const isNotFoundInOptions = $derived(value && !$options.map((o) => o.value).includes(value));
</script>

<div class="flex items-center gap-1">
	<Command.Input
		data-slot="select-input"
		class={cn(className)}
		bind:ref
		{...restProps}
		bind:value
	/>
	{#if addition && isNotFoundInOptions}
		{@render addition({
			manager: optionManager,
			accessor: {
				set input(newValue: string) {
					value = newValue;
				},
				get input() {
					return value;
				}
			}
		})}
	{/if}
</div>
