type OptionType = {
    value: any;
    label: string;
    icon?: string;
    information?: string;
    enabled?: boolean;
};

type valueSetterType = (option: OptionType) => void;
type valueGetterType = () => any;

export type {
    OptionType,
    valueSetterType,
    valueGetterType,
}