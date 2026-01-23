<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { Content, Form, setFormContext } from '@sjsf/form';

	import { type Snippet } from 'svelte';

	interface Props {
		form: FormState<Record<string, unknown>>;
		children?: Snippet;
	}

	let { form, children }: Props = $props();

	setFormContext(form);
</script>

<Form>
	<Content />
	{@render children?.()}
</Form>
