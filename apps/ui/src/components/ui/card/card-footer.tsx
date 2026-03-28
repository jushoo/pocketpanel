import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardFooterProps extends JSX.HTMLAttributes<HTMLDivElement> {}

export const CardFooter: Component<CardFooterProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<div
			data-slot="card-footer"
			class={cn('flex items-center p-6 pt-0', local.class)}
			{...rest}
		/>
	);
};
