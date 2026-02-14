<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { Content, Form, setFormContext } from '@sjsf/form';
	import { type Snippet } from 'svelte';

	interface Props {
		form: FormState<Record<string, unknown>>;
		children?: Snippet;
		ref?: HTMLFormElement;
	}

	let { form, children, ref = $bindable() }: Props = $props();

	setFormContext(form);
</script>

<Form bind:ref>
	<Content />
	{@render children?.()}
</Form>
