import { type ReactNode } from 'react';
import { type TooltipContentProps } from 'recharts';
import { KeyDate, KeyDistance, KeyTotalDistance } from './keys';

type Payload = Record<KeyDate, string> &
  Record<KeyDistance, number> &
  Record<KeyTotalDistance, number>;

interface ValuePoint {
  readonly value: number;
  readonly payload: Payload;
}

const isValuePoint = (x: unknown): x is ValuePoint =>
  x !== undefined &&
  x !== null &&
  typeof x === 'object' &&
  'value' in x &&
  typeof x.value === 'number';

export const TooltipContent = (
  props: TooltipContentProps<string | number, string>,
): ReactNode => {
  const { active, payload, label } = props;
  if (!active || payload.length === 0) {
    return null;
  }
  const valuePoint = payload.find(isValuePoint);
  if (!valuePoint) {
    return null;
  }
  return (
    <div className='card card-xs w-24 bg-base-100 shadow-sm'>
      <div className='card-body'>
        <h2 className='card-title'>{valuePoint.value}km</h2>
        <p>
          <time>{label}</time>
        </p>
        <p>総走行距離: {valuePoint.payload.総走行距離}km</p>
      </div>
    </div>
  );
};
