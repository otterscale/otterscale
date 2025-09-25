<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';

	import { INPUT_CLASSNAME } from './utils.svelte';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Switch from '$lib/components/ui/switch';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		value: checked = $bindable(undefined),
		required,
		nullable = false,
		descriptor,
		format = 'checkbox',
		invalid = $bindable(),

		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> & {
		nullable?: boolean;
		descriptor?: (v: any) => string;
		format?: 'switch' | 'checkbox';
		invalid?: boolean | null | undefined;
	} = $props();

	let proxyChecked = $state(false);

	const isInvalid = $derived(required && [null, undefined].includes(checked));
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div
	class={cn(
		INPUT_CLASSNAME,
		'relative flex items-center gap-2',
		format === 'switch' ? 'ring-1' : 'border-none shadow-none',
		isInvalid && format === 'switch' ? 'ring-destructive' : '',
		isInvalid && format === 'checkbox' ? 'ring-destructive ring-1' : '',
		className,
	)}
>
	<span class="pr-15">
		{#if required}
			{@const isValid = [true, false].includes(checked)}
			{@const isNull = [null, undefined].includes(checked)}

			{#if isValid}
				{#if descriptor}
					<p class={cn('text-sm', checked ? 'text-primary' : 'text-muted-foreground ')}>
						{descriptor(checked)}
					</p>
				{:else if checked === true}
					<Badge variant="default">True</Badge>
				{:else if checked === false}
					<Badge variant="outline">False</Badge>
				{/if}
			{:else if isNull}
				<p class="text-destructive/60 text-xs">Required</p>
			{:else}
				<Badge variant="destructive">Invalid</Badge>
			{/if}
		{:else}
			{@const isValid = [null, undefined, true, false].includes(checked)}

			{#if isValid}
				{#if descriptor}
					<p class={cn('text-sm', checked ? 'text-primary font-bold' : 'text-muted-foreground ')}>
						{descriptor(checked)}
					</p>
				{:else if checked === true}
					<Badge variant="default">True</Badge>
				{:else if checked === false}
					<Badge variant="outline">False</Badge>
				{:else if checked === null || checked === undefined}
					<Badge variant="secondary">Null</Badge>
				{/if}
			{:else}
				<Badge variant="destructive">Invalid</Badge>
			{/if}
		{/if}
	</span>

	<span
		class={cn(
			'absolute top-1/2 right-0 flex -translate-y-1/2 items-center rounded-full hover:cursor-pointer focus:outline-none',
			// format === 'checkbox' && isInvalid ? 'ring-destructive ring-1' : ''
		)}
	>
		{#if nullable}
			<button
				class="absolute top-1/2 right-6 flex -translate-y-1/2 items-center hover:cursor-pointer focus:outline-none"
				onclick={() => {
					checked = undefined;
				}}
			>
				<Icon icon="ph:x" />
			</button>
		{/if}
		{#if checked === undefined}
			<Switch.Root
				bind:ref
				bind:checked={proxyChecked}
				data-slot="input-boolean"
				{...restProps}
				onCheckedChange={() => {
					checked = proxyChecked;
				}}
			/>
		{:else}
			<Switch.Root bind:ref bind:checked data-slot="input-boolean" {...restProps} />
		{/if}
	</span>
</div>
