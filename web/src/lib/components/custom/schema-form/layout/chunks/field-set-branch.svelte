<script lang="ts">
	import { type ComponentProps, getFormContext, uiOptionProps } from '@sjsf/form';
	import { getThemeContext } from '@sjsf/shadcn4-theme';
	import type { HTMLFieldsetAttributes } from 'svelte/elements';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, config, attributes }: Props = $props();

	const ctx = getFormContext();
	const themeCtx = getThemeContext();
	const { FieldSet } = $derived(themeCtx.components);
</script>

<FieldSet {...uiOptionProps('shadcn4FieldSet')(attributes as HTMLFieldsetAttributes, config, ctx)}>
	{@render children()}
</FieldSet>
