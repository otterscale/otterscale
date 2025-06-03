<script lang="ts">
	import { BORDER_INPUT_CLASSNAME, UNFOCUS_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';
	import type { ZodFirstPartySchemaTypes } from 'zod';
	import { InputValidator } from './utils.svelte';

	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type'> & {
			type?: Exclude<HTMLInputTypeAttribute, 'file' | 'password'>;
		}
	>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		type,
		required,
		schema,
		class: className,
		...restProps
	}: Props & { schema?: ZodFirstPartySchemaTypes } = $props();
</script>

{#if schema}
	{@const validator = new InputValidator(schema)}
	{@const validation = validator.validate(value)}
	{@const requirementError = required && !value}
	{@const inputError = value && !validation.valid}
	{@const controllerClassName = requirementError || inputError ? 'ring-destructive ring-1' : ''}

	{@render Controller(controllerClassName)}
	<div class="transition-all duration-500">
		{#if requirementError}
			<div class="animate-in fade-in flex items-center gap-1">
				<Icon icon="ph:asterisk" class="text-destructive size-2" />
				<p class="text-destructive text-xs">Required</p>
			</div>
		{/if}
		{#if inputError}
			<div class="animate-in fade-in flex items-center gap-2">
				{#each validation.errors as error}
					<span class="flex items-center gap-1">
						<Icon icon="ph:warning" class="text-destructive" />
						<p class="text-destructive text-xs">{error.message}</p>
					</span>
				{/each}
			</div>
		{/if}
	</div>
{:else}
	{@render Controller()}
{/if}

{#snippet Controller(controllerClassName: string = '')}
	<div class={cn(BORDER_INPUT_CLASSNAME, controllerClassName, className)}>
		{#if type}
			<span class="pl-3">
				<Icon icon={typeToIcon[type]} />
			</span>
		{/if}
		<Input
			bind:ref
			data-slot="input-general"
			class={cn(UNFOCUS_INPUT_CLASSNAME)}
			{type}
			bind:value
			{...restProps}
		/>
	</div>
{/snippet}
