<script lang="ts">
	import { type ComponentProps, getFormContext, layoutAttributes } from '@sjsf/form';
	import { getContext, setContext } from 'svelte';

	import ArrayItemControls from './chunks/array-item-controls.svelte';
	import DefaultBranch from './chunks/default-branch.svelte';
	import FieldBranch from './chunks/field-branch.svelte';
	import FieldGroupBranch from './chunks/field-group-branch.svelte';
	import FieldLegendBranch from './chunks/field-legend-branch.svelte';
	import FieldSetBranch from './chunks/field-set-branch.svelte';
	import FieldTitleRow from './chunks/field-title-row.svelte';
	import SimpleContent from './chunks/simple-content.svelte';
	import LazyBranch from './lazy-branch.svelte';

	let props: ComponentProps['layout'] = $props();
	const { type, config } = props;

	// Use context to track recursion depth implicitly
	const parentDepth = getContext<number>('layout-depth') || 0;
	setContext('layout-depth', parentDepth + 1);

	// Threshold for switching to async rendering to break the call stack
	// First 10 layers render synchronously for performance, deeper layers use LazyBranch
	const LAZY_THRESHOLD = 20;
	const useLazy = parentDepth > LAZY_THRESHOLD;

	const ctx = getFormContext();
	const attributes = $derived(layoutAttributes(ctx, config, 'layout', 'layouts', type, {}));

	const isMeta = $derived(
		type === 'field-meta' || type === 'array-field-meta' || type === 'object-field-meta'
	);

	const BRANCHES = {
		simple: SimpleContent,
		'array-item-controls': ArrayItemControls,
		fieldset: FieldSetBranch,
		field: FieldBranch,
		'field-title-row': FieldTitleRow,
		'field-legend': FieldLegendBranch,
		'field-group': FieldGroupBranch,
		default: DefaultBranch
	};

	function getBranchType(t: string, meta: boolean, attrs: any): keyof typeof BRANCHES {
		if ((t === 'field-content' || meta) && Object.keys(attrs).length < 2) {
			return 'simple';
		}

		const mapping: Record<string, keyof typeof BRANCHES> = {
			'array-item-controls': 'array-item-controls',
			'array-field': 'fieldset',
			'object-field': 'fieldset',
			field: 'field',
			'field-title-row': 'field-title-row',
			'array-field-title-row': 'field-legend',
			'object-field-title-row': 'field-legend',
			'array-items': 'field-group',
			'object-properties': 'field-group',
			'multi-field': 'field-group',
			'multi-field-content': 'field-group'
		};

		return mapping[t] || 'default';
	}

	const branchKey = $derived(getBranchType(type, isMeta, attributes));
	const Branch = $derived(BRANCHES[branchKey]);

	// const branchProps = $derived(branchKey === 'simple' ? props : { ...props, attributes });
	const branchProps = $derived({ ...props, attributes });
</script>

{#if parentDepth > 100}
	{console.warn('Recursion Limit Reached!', parentDepth, type)}
	<div style="border: 2px solid red; padding: 10px;">
		<strong>🛑 Recursion Limit Reached ({parentDepth})!</strong>
		<p>Type: {type}</p>
		<pre>{JSON.stringify(config, null, 2)}</pre>
	</div>
{:else if useLazy}
	<!-- Use LazyBranch to break the call stack for deep recursion -->
	<LazyBranch component={Branch} {...branchProps} />
{:else}
	<Branch {...branchProps} />
{/if}
