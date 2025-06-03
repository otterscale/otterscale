type OptionType = {
    value: string;
    label: string;
    icon?: string;
    subOptions?: OptionType[];
};

type AncestralOptionType = OptionType[];

type valuesSetterType = (ancestralOptions: AncestralOptionType[]) => void


export type {
    OptionType,
    AncestralOptionType,
    valuesSetterType
}