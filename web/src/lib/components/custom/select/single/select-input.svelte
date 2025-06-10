<script lang="ts">
	import * as Command from '$lib/components/ui/command';
	import { Command as CommandPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import type { Writable } from 'svelte/store';
	import type { OptionType } from './types';
	import { getContext, type Snippet } from 'svelte';

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
