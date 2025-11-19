import type { triggerVariants } from './utils';

type TriggerVariant = VariantProps<typeof triggerVariants>['variant'];
type TriggerSize = VariantProps<typeof triggerVariants>['size'];

type Booleanified<T> = {
	[K in keyof T]: boolean;
};

export type { TriggerSize, TriggerVariant, Booleanified };
