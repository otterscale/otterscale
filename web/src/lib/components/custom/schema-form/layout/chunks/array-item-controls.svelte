<script lang="ts">
	import { type ComponentProps, getFormContext, uiOptionProps } from '@sjsf/form';
	import { getThemeContext } from '@sjsf/shadcn4-theme';

	// We extend the standard layout props to include the calculated attributes
	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, config, attributes }: Props = $props();

	const ctx = getFormContext();
	const themeCtx = getThemeContext();
	const { ButtonGroup } = $derived(themeCtx.components);
</script>

<ButtonGroup {...uiOptionProps('shadcn4ButtonGroup')(attributes, config, ctx)}>
	{@render children()}
</ButtonGroup>
