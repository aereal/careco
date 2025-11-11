import { type FC } from 'react';
import { RecordDialog } from './component.dialog';
import { OpenButton } from './component.open-button';

export const RecordDialogContainer: FC = () => {
  return (
    <>
      <OpenButton />
      <RecordDialog />
    </>
  );
};
