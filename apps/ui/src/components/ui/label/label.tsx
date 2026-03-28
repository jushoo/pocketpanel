import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface LabelProps extends JSX.LabelHTMLAttributes<HTMLLabelElement> {}

export const Label: Component<LabelProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<label
			data-slot="label"
			class={cn(
				'flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50',
				local.class
			)}
			{...rest}
		/>
	);
};
