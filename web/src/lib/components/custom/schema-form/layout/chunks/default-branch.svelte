<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { type, children, attributes }: Props = $props();

	const isItem = $derived(type === 'array-item');
	const isGrowable = $derived(
		type === 'array-item-content' ||
			type === 'object-property-key-input' ||
			type === 'object-property-content'
	);
	const isObjectProperty = $derived(type === 'object-property');
	const isMultiFieldControls = $derived(type === 'multi-field-controls');
</script>

<div
	class={{
		grow: isGrowable,
		'flex items-center gap-2': isMultiFieldControls,
		'flex items-start gap-1.5': isItem,
		'grid grid-cols-1 grid-rows-[1fr] items-start gap-x-1.5 [&:has(>:nth-child(2))]:grid-cols-[1fr_1fr_auto]':
			isObjectProperty
	}}
	{...attributes}
>
	{@render children()}
</div>
