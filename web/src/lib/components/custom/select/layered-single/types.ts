type OptionType = {
    value: string;
    label: string;
    icon?: string;
    subOptions?: OptionType[];
};

type AncestralOptionType = OptionType[];

type valueSetterType = (newVamue: any) => void
type valueGetterType = () => any

export type {
    OptionType,
    AncestralOptionType,
    valueSetterType,
    valueGetterType
}