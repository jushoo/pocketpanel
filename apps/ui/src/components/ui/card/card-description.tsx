import { splitProps } from 'solid-js';
import type { JSX, Component } from 'solid-js';
import { cn } from '~/lib/utils';

export interface CardDescriptionProps extends JSX.HTMLAttributes<HTMLParagraphElement> {}

export const CardDescription: Component<CardDescriptionProps> = (props) => {
	const [local, rest] = splitProps(props, ['class']);

	return (
		<p
			data-slot="card-description"
			class={cn('text-sm text-muted-foreground', local.class)}
			{...rest}
		/>
	);
};
