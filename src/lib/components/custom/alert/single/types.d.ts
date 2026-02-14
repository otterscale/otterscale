import { alertVariants } from './alert.svelte';

type AlertType = {
	title: string;
	message: string;
	action?: () => void;
	variant?: alertVariants;
};

export type { AlertType };
