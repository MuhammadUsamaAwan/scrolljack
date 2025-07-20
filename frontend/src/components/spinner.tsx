import { Loader2 } from 'lucide-react';

export function Spinner() {
  return (
    <div className='flex justify-center'>
      <Loader2 className='animate-spin h-6 w-6' />
    </div>
  );
}
