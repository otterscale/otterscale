<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Switch from '$lib/components/ui/switch';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		value: checked = $bindable(undefined),
		id,
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
		'relative flex items-center gap-2 ring-1',
		isInvalid ? 'ring-destructive' : '',
		format === 'checkbox' && 'border-none shadow-none ring-0',
		className
	)}
>
	{#if format === 'switch'}
		<span class="absolute left-3 top-1/2 -translate-y-1/2 items-center">
			<Icon icon={typeToIcon['boolean']} />
		</span>
	{/if}

	<span class={cn('pr-15', format === 'switch' ? 'pl-9' : '')}>
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
			'absolute right-3 top-1/2 flex -translate-y-1/2 items-center rounded-full hover:cursor-pointer focus:outline-none',
			format === 'checkbox' && isInvalid ? 'ring-destructive ring-1' : ''
		)}
	>
		{#if nullable}
			<button
				class="absolute right-9 top-1/2 flex -translate-y-1/2 items-center hover:cursor-pointer focus:outline-none"
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
