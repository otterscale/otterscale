<script lang="ts" module>
	const type = 'color';
</script>

<script lang="ts">
	import { BORDER_INPUT_CLASSNAME, UNFOCUS_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';

	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input';
	import { z, type ZodFirstPartySchemaTypes } from 'zod';
	import { InputValidator } from './utils.svelte';

	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from 'svelte/elements';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type'> & {
			type?: Exclude<HTMLInputTypeAttribute, 'file' | 'password'>;
		}
	>;

	let {
		ref = $bindable(null),
		value = $bindable(),
		schema = z.string().regex(/^#[0-9a-fA-F]{6}$/),
		class: className,
		...restProps
	}: Props & {
		schema?: ZodFirstPartySchemaTypes;
	} = $props();
</script>

{#if schema}
	{@const validator = new InputValidator(schema)}
	{@const validation = validator.validate(value)}
	{@const inputError = value && !validation.valid}
	{@const controllerClassName = inputError ? 'ring-destructive ring-1' : ''}

	<div class={cn(BORDER_INPUT_CLASSNAME, controllerClassName, 'h-10 justify-between', className)}>
		<span class="flex items-center gap-2">
			<span class="flex items-center gap-2">
				<span class="pl-3">
					<Icon icon={typeToIcon[type]} />
				</span>
				<Badge variant="outline">{value}</Badge>
			</span>
		</span>
		<Input
			bind:ref
			data-slot="input-color"
			class={cn(UNFOCUS_INPUT_CLASSNAME, 'mr-3 aspect-square h-7 w-fit cursor-pointer p-0')}
			{type}
			bind:value
			{...restProps}
		/>
	</div>
	<div class="transition-all duration-500">
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
{/if}
