<script lang="ts">
	import {
		type ComponentProps,
		getComponent,
		getFieldComponent,
		getFieldErrors,
		getFormContext,
		Text
	} from '@sjsf/form';
	import { getArrayContext } from '@sjsf/form/fields/array/context.svelte';

	import Template from './array-template.svelte';

	let {
		index,
		value = $bindable(),
		config,
		uiOption,
		translate
	}: ComponentProps['arrayItemField'] = $props();

	const ctx = getFormContext();
	const arrayCtx = getArrayContext();

	const Field = $derived(getFieldComponent(ctx, config));
	const Button = $derived(getComponent(ctx, 'button', config));

	const canCopy = $derived(arrayCtx.canCopy(index));
	const canRemove = $derived(arrayCtx.canRemove(index));
	const canMoveUp = $derived(arrayCtx.canMoveUp(index));
	const canMoveDown = $derived(arrayCtx.canMoveDown(index));
	const toolbar = $derived(canCopy || canRemove || canMoveUp || canMoveDown);
	const errors = $derived(getFieldErrors(ctx, config.path));
</script>

{#snippet buttons()}
	<div
		class="flex items-center gap-2 [&_button]:size-7 [&_button]:border-none [&_button]:bg-transparent [&_button]:shadow-none"
	>
		{#if arrayCtx.orderable()}
			<Button
				{errors}
				{config}
				type="array-item-move-up"
				disabled={!canMoveUp}
				onclick={() => {
					arrayCtx.moveItemUp(index);
				}}
			>
				<Text {config} id="move-array-item-up" {translate} />
			</Button>
			<Button
				{errors}
				{config}
				disabled={!canMoveDown}
				type="array-item-move-down"
				onclick={() => {
					arrayCtx.moveItemDown(index);
				}}
			>
				<Text {config} id="move-array-item-down" {translate} />
			</Button>
		{/if}
		{#if canCopy}
			<Button
				{errors}
				{config}
				type="array-item-copy"
				onclick={() => {
					arrayCtx.copyItem(index);
				}}
				disabled={false}
			>
				<Text {config} id="copy-array-item" {translate} />
			</Button>
		{/if}
		{#if canRemove}
			<div class="[&_button]:*:text-destructive">
				<Button
					{errors}
					{config}
					disabled={false}
					type="array-item-remove"
					onclick={() => {
						arrayCtx.removeItem(index);
					}}
				>
					<Text {config} id="remove-array-item" {translate} />
				</Button>
			</div>
		{/if}
	</div>
{/snippet}
<div class="border-l-2 py-4 pl-6">
	<Template
		type="template"
		{index}
		{value}
		{config}
		{errors}
		buttons={toolbar ? buttons : undefined}
		{uiOption}
	>
		<Field type="field" bind:value={value as undefined} {config} {uiOption} {translate} />
	</Template>
</div>
