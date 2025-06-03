type OptionType = {
    value: string;
    label: string;
    icon?: string;
    subOptions?: OptionType[];
};

type AncestralOptionType = OptionType[];

type valueSetterType = (ancestralOption: AncestralOptionType) => void

export type {
    OptionType,
    AncestralOptionType,
    valueSetterType
}