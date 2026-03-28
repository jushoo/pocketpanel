import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardProps extends JSX.HTMLAttributes<HTMLDivElement> {}

export const Card: Component<CardProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<div
			data-slot="card"
			class={cn(
				'rounded-xl border bg-card text-card-foreground shadow',
				local.class
			)}
			{...rest}
		/>
	);
};
