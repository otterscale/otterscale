<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';
	import { getThemeContext } from '@sjsf/shadcn4-theme';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, attributes }: Props = $props();

	const themeCtx = getThemeContext();
	const { FieldLegend } = $derived(themeCtx.components);
</script>

<FieldLegend class="flex w-full items-center justify-between" {...attributes}>
	{@render children()}
</FieldLegend>
