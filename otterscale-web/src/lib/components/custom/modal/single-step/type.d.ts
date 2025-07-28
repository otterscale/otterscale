import type { triggerVariants } from "./utils";

type TriggerVariant = VariantProps<typeof triggerVariants>["variant"]
type TriggerSize = VariantProps<typeof triggerVariants>["size"];

export type { TriggerVariant, TriggerSize };
