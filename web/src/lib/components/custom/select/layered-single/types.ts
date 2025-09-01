type OptionType = {
	value: string;
	label: string;
	icon?: string;
	subOptions?: OptionType[];
};
type AncestralOptionType = OptionType[];
type AccessorType = {
	value: any;
};

export type { AccessorType, AncestralOptionType, OptionType };
