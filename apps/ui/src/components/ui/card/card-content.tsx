import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardContentProps extends JSX.HTMLAttributes<HTMLDivElement> {}

export const CardContent: Component<CardContentProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<div
			data-slot="card-content"
			class={cn('p-6 pt-0', local.class)}
			{...rest}
		/>
	);
};
