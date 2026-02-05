<script lang="ts">
	import { getFormContext, uiOptionProps, type ComponentProps } from '@sjsf/form';
	import type { HTMLFieldsetAttributes } from 'svelte/elements';
	import { getThemeContext } from '@sjsf/shadcn4-theme';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, config, attributes }: Props = $props();

	const ctx = getFormContext();
	const themeCtx = getThemeContext();
	const { FieldSet } = $derived(themeCtx.components);
</script>

<FieldSet {...uiOptionProps('shadcn4FieldSet')(attributes as HTMLFieldsetAttributes, config, ctx)}>
	{@render children()}
</FieldSet>
