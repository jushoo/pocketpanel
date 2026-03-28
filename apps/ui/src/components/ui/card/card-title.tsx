import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardTitleProps extends JSX.HTMLAttributes<HTMLHeadingElement> {}

export const CardTitle: Component<CardTitleProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<h3
			data-slot="card-title"
			class={cn('font-semibold leading-none tracking-tight', local.class)}
			{...rest}
		/>
	);
};
