type OptionType = {
    value: string;
    label: string;
    icon?: string;
    subOptions?: OptionType[];
};

type AncestralOptionType = OptionType[];

type valuesSetterType = (newValues: any[]) => void
type valuesGetterType = () => any[]


export type {
    OptionType,
    AncestralOptionType,
    valuesSetterType,
    valuesGetterType
}